package config

import (
	"fmt"
	"log/slog"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy/envfileloader"
)

var loader *envfileloader.Loader[root]

func LoadConfig(keyDebugMode string) {
	envLoader, err := envfileloader.New(func(t *root, err error) {
		if err != nil {
			slog.Error("failed to load env", "error", err)
			return
		}

		slog.Debug("env reloaded", "config", t)
	},
		envfileloader.WithCallbackOnChangeWhenOnKeyTrue(keyDebugMode),
		envfileloader.WithFiles("env.json", "../env.json", "../../env.json", "../../../env.json"),
		envfileloader.WithFileType("json"),
		envfileloader.WithWatch(true),
		envfileloader.WithTag("env"),
	)
	if err != nil {
		panic(err)
	}

	loader = envLoader
}

func GetAppRestApi() AppRestApi {
	return loader.Get().AppRestApi
}

func GetAppGrpcApi() AppGrpcApi {
	return loader.Get().AppGrpcApi
}

func GetAppScheduler() AppScheduler {
	return loader.Get().AppScheduler
}

func GetPprofAppGrpcApi() Pprof {
	return loader.Get().AppGrpcApi.Pprof
}

func GetPprofAppRestApi() Pprof {
	return loader.Get().AppRestApi.Pprof
}

func GetPprofAppScheduler() Pprof {
	return loader.Get().AppScheduler.Pprof
}

func GetDatabase() Database {
	return loader.Get().Database
}

func UnwatchLoader() error {
	slog.Info("unwatching config loader")
	err := loader.Unwatch()
	if err != nil {
		return fmt.Errorf("failed to unwatch config loader: %w", err)
	}
	return nil
}
