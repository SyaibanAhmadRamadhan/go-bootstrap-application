package provider

import (
	"context"
	"erp-directory-service/internal/config"
	"fmt"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/observability/zerologhook"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/generic"
	"github.com/rs/zerolog"
)

func NewLogging() {
	appConfig := config.GetApp()
	observability.NewLog(observability.LogConfig{
		Hook:        zerologhook.NewRotatingWriter(fmt.Sprintf("%s.log", appConfig.Name), 10, 2, 30, true),
		Mode:        "json",
		Level:       generic.Ternary(appConfig.Env == "development", "debug", "info"),
		Env:         appConfig.Env,
		ServiceName: appConfig.Name,
	})
	observability.Start(context.Background(), zerolog.InfoLevel).Msg("init logging successfully")
}
