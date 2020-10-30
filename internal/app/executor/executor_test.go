package executor_test

import (
	"Nani/internal/app/cache"
	config2 "Nani/internal/app/config"
	"Nani/internal/app/executor"
	"Nani/internal/app/inhuman"
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func exe(newstorage bool) *executor.Executor {
	os.Setenv("api_key", "Security 3923cf9a417e73be95b40dc5db60c97dcb876a61")
	config := config2.New()
	api := inhuman.New(config)
	storage := cache.New(newstorage)
	return executor.New(api, storage, config)
}

func TestScrap_ShouldStartingScrapingFor5Seconds_NoError(t *testing.T)  {
	ex := exe(true)

	go func() {
		time.Sleep(time.Second * 5)
		ex.Stop()
	}()
	assert.NoError(t, ex.Scrap(context.Background(), "../../../bundles.txt"))
}
