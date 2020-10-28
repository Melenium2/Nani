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
)

type ExecutorError struct {
	t      string
	er     string
	bundle string
}

type databaseCh chan interface{}
type errorCh chan ExecutorError

type Executor struct {
	externalApi inhuman.ExternalApi
	cache       cache.Storage
	ctx         context.Context
	repository  db.AppRepository
	keyCache    cache.KeyStorage
	config      *config.Config
	db          databaseCh
	errs        errorCh
}

func (ex *Executor) Scrap(ctx context.Context, scrapfile string) error {
	ex.ctx = ctx

	err := ex.declareTask(scrapfile)
	if err != nil {
		return err
	}
	go ex.selector()
	go ex.appsBatch()

	// If 'last' exists then slice bundles array and continue from last parsed app
	cachedBundles, err := ex.cache.GetV("bundles")
	if err != nil {
		return err
	}
	bundles, ok := cachedBundles.([]string)
	if !ok {
		return errors.New("can not convert bundles")
	}
	ex.storeApps(true, bundles...)
	// Run function for scraping application
	// Then storeApps this application info to apps table with all application bundles
	// In the same time we get top keywords from application and send them to channel
	// which get 250 apps developers and save them (or not)
	// Secluded channel will scrap info about application from previous step
	// and save them to apps table too

	// Errors...
	// If errors occur then we need to storeApps state at the cache (redis or local)
	// and keep thinking over about this question
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
			ex.saveError("apps", fmt.Sprintf("error in external api method App() %s", err), v)
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
		ex.saveError("keys", fmt.Sprintf("error in external method Keys() %s", err), app.Bundle)
	}
	ex.db <- keys
}

// saveError save error to local cache
// @params
//	t: string (type of error)
// 	er: string (error representation)
// 	bundle: string (bundle where error occurred)
func (ex *Executor) saveError(t, er, bundle string) {
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
	for {
		select {
		case t := <-ex.db:
			switch data := t.(type) {
			case *inhuman.App:
				apps = append(apps, data)
				if len(apps) > 50 {
					err := ex.repository.InsertBatch(ex.ctx, apps)
					if err != nil {
						ex.saveError("db", err.Error(), "")
						continue
					}
					apps = nil
				}
			case inhuman.Keywords:
				keys := SortKeywords(data)[:ex.config.KeysCount]
				for _, k := range keys {
					err := ex.keyCache.Set(k)
					if err != nil {
						ex.saveError("keyCache", err.Error(), "")
					}
				}
			}
		default:
			// Add case for error or cancel loop
		}
	}
}

func (ex *Executor) appsBatch() {
	for {
		key, err := ex.keyCache.Next()
		if err != nil {
			if err.Error() == "keywords are out of range" {
				return
			}
			continue
		}
		ex.storeApps(false, key)
	}
}

func (ex *Executor) Stop() {

}

func New(api inhuman.ExternalApi, storage cache.Storage) *Executor {
	return &Executor{
		externalApi: api,
		cache:       storage,
	}
}
