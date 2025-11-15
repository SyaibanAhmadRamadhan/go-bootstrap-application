package app

import (
	"context"
	"erp-directory-service/gen/restapigen"
	"erp-directory-service/internal/config"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/chix"
)

type restapi struct {
	server *http.Server
	port   int
}

func (r *restapi) Shutdown(ctx context.Context) error {
	err := r.server.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *restapi) Start() {
	fmt.Println("REST API listening on", r.server.Addr)
	err := r.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
	}
}

func NewRestApi(port int) *restapi {
	handler := chix.New(chix.Config{
		BlacklistRouteLogResponse: map[string]struct{}{},
		SensitiveFields:           map[string]struct{}{},
		CorsConf: chix.CorsConfig{
			AllowOrigins:     nil,
			AllowMethods:     nil,
			AllowHeaders:     nil,
			AllowCredentials: true,
		},
		AppName: config.GetApp().Name,
		UseOtel: false,
	})

	router := router{}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: restapigen.HandlerFromMux(router, handler),
	}

	return &restapi{
		server: srv,
		port:   port,
	}
}
