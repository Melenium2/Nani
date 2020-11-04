package main

import (
	"Nani/internal/app/cache"
	"Nani/internal/app/config"
	"Nani/internal/app/executor"
	"Nani/internal/app/inhuman"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var bundles string
	bundles = os.Getenv("bundle")

	var configDir string
	flag.StringVar(&configDir, "config", "config/dev.yml", "Application config file")
	var schemaDir string
	flag.StringVar(&schemaDir, "schema", "config/schema.sql", "Database schema")
	var cacheDir string
	flag.StringVar(&cacheDir, "cache", "internal/app/cache/cache.json", "Cache file dir")
	if bundles == "" {
		flag.StringVar(&bundles, "e", "bundles.txt", "The file from which to parse")
	}
	flag.Parse()


	// Just for test case
	//os.Setenv("api_key", "Security 3923cf9a417e73be95b40dc5db60c97dcb876a61")
	conf := config.New(configDir)
	conf.Database.Schema = schemaDir
	conf.KeysCount = 10
	conf.AppsCount = 250

	api := inhuman.New(conf)

	storage := cache.New(false, cacheDir)

	ex := executor.New(api, storage, conf)

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGINT)
		<-sig
		ex.Stop()
	}()

	err := ex.Scrap(context.Background(), bundles)
	if err != nil {
		log.Fatal(err)
	}
}
