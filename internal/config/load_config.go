package config

import (
	"fmt"
	"log/slog"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy/envfileloader"
)

var loader *envfileloader.Loader[root]
var cmdName string

func LoadConfig(cmd string) {
	keyDebugMode := ""
	switch cmd {
	case "scheduler":
		keyDebugMode = "app_scheduler.debug_mode"
	case "restapi":
		keyDebugMode = "app_rest_api.debug_mode"
	case "grpcapi":
		keyDebugMode = "app_rest_api.debug_mode"
	default:
		panic("unknown cmd name")
	}

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
	cmdName = cmd
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

func GetPprof() Pprof {
	switch cmdName {
	case "scheduler":
		return loader.Get().AppScheduler.Pprof
	case "restapi":
		return loader.Get().AppRestApi.Pprof
	case "grpcapi":
		return loader.Get().AppGrpcApi.Pprof
	default:
		slog.Error("unknown cmd name for get pprof config")
		return Pprof{}
	}
}

func GetDatabase() Database {
	switch cmdName {
	case "scheduler":
		return loader.Get().AppScheduler.Database
	case "restapi":
		return loader.Get().AppRestApi.Database
	case "grpcapi":
		return loader.Get().AppGrpcApi.Database
	default:
		slog.Error("unknown cmd name for get database config")
		return Database{}
	}
}

func GetDebugMode() bool {
	switch cmdName {
	case "scheduler":
		return loader.Get().AppScheduler.DebugMode
	case "restapi":
		return loader.Get().AppRestApi.DebugMode
	case "grpcapi":
		return loader.Get().AppGrpcApi.DebugMode
	default:
		slog.Error("unknown cmd name for get app debug mode")
		return false
	}
}

func GetAppName() string {
	switch cmdName {
	case "scheduler":
		return loader.Get().AppScheduler.Name
	case "restapi":
		return loader.Get().AppRestApi.Name
	case "grpcapi":
		return loader.Get().AppGrpcApi.Name
	default:
		slog.Error("unknown cmd name for get app name")
		return "unknown"
	}
}

func GetEnv() string {
	switch cmdName {
	case "scheduler":
		return loader.Get().AppScheduler.Env
	case "restapi":
		return loader.Get().AppRestApi.Env
	case "grpcapi":
		return loader.Get().AppGrpcApi.Env
	default:
		slog.Error("unknown cmd name for get app env")
		return "unknown"
	}
}

func UnwatchLoader() error {
	slog.Info("unwatching config loader")
	err := loader.Unwatch()
	if err != nil {
		return fmt.Errorf("failed to unwatch config loader: %w", err)
	}
	return nil
}
