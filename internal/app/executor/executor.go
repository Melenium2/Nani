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
	"time"
)

type Executor struct {
	externalApi inhuman.ExternalApi
	cache       cache.Storage
	ctx         context.Context
	repository  db.AppRepository
	keyCache    cache.KeyStorage
	config      *config.Config
	db          databaseCh
	wait        chan struct{}
	cancel      bool
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

	last, err := ex.cache.GetV("last")
	startAt := 0
	if err == nil {
		startAt = last.(int)
	}
	cachedBundles, err := ex.cache.GetV("bundles")
	if err != nil {
		return err
	}
	bundles, ok := cachedBundles.([]string)
	if !ok {
		return errors.New("can not convert bundles")
	}
	ex.storeApps(true, bundles[startAt:]...)

	<-ex.wait
	return nil
}

// declareTask method loads bundles from a file and save them to cache array
// @params
//	path: String (path to file)
// @return
// 	error: Error
func (ex *Executor) declareTask(path string) error {
	f := file.New(path)
	bundles, err := f.ReadAllSlice("\n")
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
		app, err := ex.externalApi.App(v)
		if err != nil {
			ex.saveError("apps", v, fmt.Errorf("error in external api method App() %s", err))
		} else {
			ex.db <- app
			if withKeys {
				go ex.storeKeywords(app)
			}
		}
		ex.cache.Set("last", bundles[i])
	}
}

// storeKeywords method fetching top keywords from external api and store their to db
// @params
//	app: *inhuman.App (application)
func (ex *Executor) storeKeywords(app *inhuman.App) {
	keys, err := ex.externalApi.Keys(app.Title, app.Description, app.ShortDescription, "")
	if err != nil {
		ex.saveError("keys", app.Bundle, fmt.Errorf("error in external method Keys() %s", err))
	}
	ex.db <- keys
}

// saveError save error to local cache
// @params
//	t: string (type of error)
// 	er: string (error representation)
// 	bundle: string (bundle where error occurred)
func (ex *Executor) saveError(t, bundle string, er error) {
	e, err := ex.cache.GetV("errors")
	if err != nil {
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
	for ex.cancel {
		select {
		case t := <-ex.db:
			switch data := t.(type) {
			case *inhuman.App:
				apps = append(apps, data)
				if len(apps) > 50 {
					err := ex.repository.InsertBatch(ex.ctx, apps)
					if err != nil {
						ex.saveError("db", "", err)
						continue
					}
					apps = nil
				}
			case inhuman.Keywords:
				keys := SortKeywords(data)[:ex.config.KeysCount]
				for _, k := range keys {
					err := ex.keyCache.Set(k)
					if err != nil {
						ex.saveError("keyCache", "", err)
					}
				}
			}
		default:
		}
	}

	for _, app := range apps {
		err := ex.repository.Insert(ex.ctx, app)
		if err != nil {
			ex.saveError("db", "", err)
			continue
		}
	}
}

// AppsBatch scrap new applications while keys still remain
func (ex *Executor) appsBatch() {
	for ex.cancel {
		key, err := ex.keyCache.Next()
		if err != nil {
			if err.Error() == "keywords are out of range" {
				ex.Stop()
				return
			}
			continue
		}
		ex.storeApps(false, key)
	}
}

// Stop scraping
func (ex *Executor) Stop() {
	ctx, cancel := context.WithTimeout(ex.ctx, time.Second*15)
	defer cancel()
	ex.cancel = true
	ex.ctx = ctx
	ex.wait <- struct{}{}
}

// Create new instance of Executor
func New(api inhuman.ExternalApi, storage cache.Storage) *Executor {
	return &Executor{
		externalApi: api,
		cache:       storage,
	}
}
