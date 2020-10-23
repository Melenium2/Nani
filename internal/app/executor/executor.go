package executor

import (
	"Nani/internal/app/inhuman"
	"context"
)

type Executor struct {
	externalApi inhuman.ExternalApi
	close       context.CancelFunc
}

func (ex *Executor) Scrap(ctx context.CancelFunc) {
	// Run function for scraping application
	// Then store this application info to apps table with all application bundles
	// In the same time we get top keywords from application and send them to channel
	// which get 250 apps developers and save them (or not)
	// Secluded channel will scrap info about application from previous step
	// and save them to apps table too

	// Errors...
	// If errors occur then we need to store state at the cache (redis or local)
	// and keep thinking over about this question
}

func (ex *Executor) Stop() {

}

func New(api inhuman.ExternalApi) *Executor {
	return &Executor{
		externalApi: api,
	}
}
