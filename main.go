package main

import (
	"Nani/internal/app/cache"
	"Nani/internal/app/config"
	"Nani/internal/app/executor"
	"Nani/internal/app/inhuman"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Just for test case
	os.Setenv("api_key", "Security 3923cf9a417e73be95b40dc5db60c97dcb876a61")
	conf := config.New("config/dev.yml")
	conf.Database.Schema = "config/schema.sql"
	conf.KeysCount = 10
	conf.AppsCount = 250

	api := inhuman.New(conf)

	storage := cache.New(false, "internal/app/cache/cache.json")

	ex := executor.New(api, storage, conf)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
		<-sig
		ex.Stop()
	}()

	err := ex.Scrap(context.Background(), "bundles.txt")
	if err != nil {
		log.Fatal(err)
	}
}
