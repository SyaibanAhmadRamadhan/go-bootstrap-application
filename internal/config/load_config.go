package config

import (
	"log/slog"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy/envfileloader"
)

var loader *envfileloader.Loader[root]

func LoadConfig() {
	envLoader, err := envfileloader.New(func(t *root, err error) {
		if err != nil {
			slog.Error("failed to load env", "error", err)
			return
		}

		slog.Debug("env reloaded", "config", t)
	},
		envfileloader.WithCallbackOnChangeWhenOnKeyTrue("debug_mode"),
		envfileloader.WithFiles("env.json", "../env.json", "../../env.json", "../../../env.json"),
		envfileloader.WithFileType("json"),
	)
	if err != nil {
		panic(err)
	}

	loader = envLoader
}

func GetApp() App {
	return loader.Get().App
}
