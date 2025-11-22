package infrastructure

import (
	"context"
	"fmt"
	"go-bootstrap/internal/config"
	"io"
	"os"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability/loghook"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/generic"
	"github.com/rs/zerolog"
)

func NewLogging(slogHookOption, zerologHookOption string) func() error {
	closeFn := make([]func(), 0, 2)

	var slogHook io.Writer
	switch slogHookOption {
	case "file-writer":
		rotatingWriter := loghook.NewRotatingWriter(fmt.Sprintf("%s-secondary.log", config.GetAppName()), 10, 2, 30, true)
		slogHook = rotatingWriter
		closeFn = append(closeFn, rotatingWriter.Close)
	case "std-out":
		slogHook = os.Stdout
	default:
		panic("unknown slog handler option")
	}

	var zerologHook io.Writer
	switch zerologHookOption {
	case "file-writer":
		rotatingWriter := loghook.NewRotatingWriter(fmt.Sprintf("%s-primary.log", config.GetAppName()), 10, 2, 30, true)
		zerologHook = rotatingWriter
		closeFn = append(closeFn, rotatingWriter.Close)
	case "std-out":
		zerologHook = os.Stdout
	default:
		panic("unknown slog handler option")
	}

	observability.NewLog(observability.LogConfig{
		ZerologHook: zerologHook,
		SlogHook:    slogHook,
		Mode:        "json",
		Level:       generic.Ternary(config.GetEnv() == "development", "debug", "info"),
		Env:         config.GetEnv(),
		ServiceName: config.GetAppName(),
	})
	observability.Start(context.Background(), zerolog.InfoLevel).Msg("init logging successfully")

	return func() error {
		for _, v := range closeFn {
			v()
		}
		return nil
	}
}
