package executor

import (
	"Nani/internal/app/cache"
	"Nani/internal/app/config"
	"Nani/internal/app/db"
	"Nani/internal/app/file"
	"Nani/internal/app/inhuman"
	"context"
	"errors"
	"fmt"
	murlog "github.com/Melenium2/Murlog"
	"log"
	"math"
	"runtime"
	"time"
)

/*
	TODO
	*Записывать гео по приложению
	*Дописать вторую часть процесса, в которой мы прогоняем
		всех девелоперов который спарсили и берем их приложения
	*Логгировать процесс

	*Тесты для дб. Что будет если попробывать сгенерить схему несколько раз
	*Скипать интеграционные тесты если есть на то причины

	*Ошибки сохраняются, но без самой ошибки (сделать стринг вместо ерор)
	*Ошибка которая множит слова скорее всего зависит от метода который сортирует мапу
		ошибка возникает потому что в мапе присутствуют одинаковые веса, поэтому есть
		смысл сортировать мапу подругому или делать это на беке
*/

type Executor struct {
	externalApi inhuman.ExternalApi
	cache       cache.Storage
	ctx         context.Context
	repository  db.AppRepository
	keyCache    cache.KeyStorage
	config      config.Config
	db          databaseCh
	wait        chan struct{}
	cancel      bool
	logger      murlog.Logger
}

// Scrap starting scraping all apps from scrapfile until error or
// cancel of scraping
func (ex *Executor) Scrap(ctx context.Context, scrapfile string) error {
	ex.ctx = ctx

	err := ex.declareTask(scrapfile)
	if err != nil {
		return err
	}
	go ex.selector()
	go ex.appsBatch()

	last, _ := ex.cache.GetV("last")
	cachedBundles, err := ex.cache.GetV("bundles")
	if err != nil {
		return err
	}
	bundles, ok := cachedBundles.([]string)
	startAt := 0
	if last != nil {
		l := last.(string)
		for i, b := range bundles {
			if b == l {
				startAt = i
			}
		}
	}
	if !ok {
		return errors.New("can not convert bundles")
	}
	ex.storeApps(true, bundles[startAt:]...)

	<-ex.wait

	// Second step
	return nil
}

// declareTask method loads bundles from a file and save them to cache array
// @params
//	path: String (path to file)
// @return
// 	error: Error
func (ex *Executor) declareTask(path string) error {
	f := file.New(path)
	sep := "\n"
	if runtime.GOOS == "windows" {
		sep = "\r\n"
	}
	bundles, err := f.ReadAllSlice(sep)
	if err != nil {
		return err
	}
	ex.cache.Set("bundles", bundles)
	return nil
}

// storeApps method downloading information about given slice of bundles and store their to db
// and if withKeys flag is true, then starting fetching keywords from this bundles
// @params
//	withKeys: bool (store keywords from app or not)
//	bundles: ...string (bundles for scraping)
func (ex *Executor) storeApps(withKeys bool, bundles ...string) {
	for i, v := range bundles {
		if ex.cancel {
			break
		}

		app, err := ex.externalApi.App(v)
		if err != nil {
			ex.logger.Log("log", err, "Bundle", v)
			ex.saveError("apps", v, fmt.Errorf("error in external api method App() %s", err))
		} else {
			ex.db <- app
			if withKeys {
				go ex.storeKeywords(app)
				ex.cache.Set("last", bundles[i])
			}
		}
	}
}

// storeKeywords method fetching top keywords from external api and store their to db
// @params
//	app: *inhuman.App (application)
func (ex *Executor) storeKeywords(app *inhuman.App) {
	if ex.cancel {
		return
	}

	keys, err := ex.externalApi.Keys(app.Title, app.Description, app.ShortDescription, "")
	if err != nil {
		ex.logger.Log("log", err)
		ex.saveError("keys", app.Bundle, fmt.Errorf("error in external method Keys() %s", err))
	}
	ex.db <- keys
}

// saveError save error to local cache
// @params
//	T: string (type of error)
// 	Er: string (error representation)
// 	Bundle: string (Bundle where error occurred)
func (ex *Executor) saveError(t, bundle string, er error) {
	e, err := ex.cache.GetV("errors")
	if err != nil {
		ex.logger.Log("log", err)
		ex.cache.Set("errors", []ExecutorError{{t, er, bundle}})
		return
	}
	appErrors := e.([]ExecutorError)
	appErrors = append(appErrors, ExecutorError{t, er, bundle})
	ex.cache.Set("errors", appErrors)
}

// selector main loop of channels
func (ex *Executor) selector() {
	apps := make([]*inhuman.App, 0)

	for t := range ex.db {
		switch data := t.(type) {
		case *inhuman.App:
			apps = append(apps, data)
			if len(apps) > 50 {
				err := ex.repository.InsertBatch(ex.ctx, apps)
				if err != nil {
					ex.logger.Log("log", err)
					ex.saveError("Db", "", err)
					continue
				}
				apps = nil
			}
		case inhuman.Keywords:
			s := int(math.Min(float64(ex.config.KeysCount), float64(len(data))))
			log.Print(data)
			keys := SortKeywords(data)[:s]
			log.Printf("%v", keys)
			for _, k := range keys {
				err := ex.keyCache.Set(k)
				if err != nil {
					ex.logger.Log("log", err)
					ex.saveError("keyCache", "", err)
				}
			}
		}
	}

	for _, app := range apps {
		err := ex.repository.Insert(ex.ctx, app)
		if err != nil {
			ex.logger.Log("log", err)
			ex.saveError("Db", "", err)
			continue
		}
	}
}

// AppsBatch scrap new applications while keys still remain
func (ex *Executor) appsBatch() {
	for !ex.cancel {
		key, err := ex.keyCache.Next()
		if err != nil {
			if err.Error() == "keywords are out of range" {
				break
			}
			continue
		}
		ex.logger.Log("next key", key)
		res, err := ex.externalApi.Flow(key)
		if err != nil {
			ex.logger.Log("log", err)
			ex.saveError("keys", key, err)
			ex.keyCache.Rollback()
			continue
		}
		bundles := make([]string, len(res))
		for i := 0; i < len(bundles); i++ {
			bundles[i] = res[i].Bundle
		}
		ex.storeApps(false, bundles...)
	}
	ex.wait <- struct{}{}
}

// Stop scraping
func (ex *Executor) Stop() {
	//ctx, cancel := context.WithTimeout(ex.ctx, time.Second*15)
	//ex.ctx = ctx
	//defer cancel()
	ex.logger.Log("log", "Starting stopping application")
	ex.cancel = true
	<-ex.wait
	ex.wait <- struct{}{}
	time.Sleep(time.Second * 3)
	close(ex.db)
	ex.logger.Log("log", "Closing")
}

// Create new instance of Executor
func New(api inhuman.ExternalApi, storage cache.Storage, config config.Config) *Executor {
	mConfig := murlog.NewConfig()
	mConfig.CallerPref()
	mConfig.TimePref(time.RFC1123)

	return &Executor{
		externalApi: api,
		cache:       storage,
		keyCache:    cache.NewKeyCache(storage),
		repository:  db.New(config.Database),
		config:      config,
		db:          make(databaseCh, 5),
		wait:        make(chan struct{}, 1),
		cancel:      false,
		logger:      murlog.NewLogger(mConfig),
	}
}
