package app

import (
	"context"
	"erp-directory-service/internal/config"
	"erp-directory-service/internal/gen/restapigen"
	healthcheckrepository "erp-directory-service/internal/module/healthcheck/repository"
	healthcheckservice "erp-directory-service/internal/module/healthcheck/service"
	"erp-directory-service/internal/provider"
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/chix"
	"github.com/go-chi/chi/v5"
)

type restapi struct {
	server  *http.Server
	port    int
	closeFn []func() error
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

	restApi := &restapi{
		port:    port,
		closeFn: make([]func() error, 0),
	}

	restApi.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: restapigen.HandlerFromMux(restApi.init(handler), handler),
	}

	return restApi
}

func (r *restapi) Shutdown(ctx context.Context) error {
	errs := make([]error, 0, len(r.closeFn))

	err := r.server.Shutdown(ctx)
	if err != nil {
		errs = append(errs, err)
	}

	for _, v := range r.closeFn {
		if err = v(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func (r *restapi) Start() {
	slog.Info(fmt.Sprintf("REST API listening on %s", r.server.Addr))
	err := r.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
	}
}

func (r *restapi) init(c *chi.Mux) routerRestApi {
	db, err := provider.NewDB()
	if err != nil {
		panic(err)
	}
	r.closeFn = append(r.closeFn, db.Close)

	healthcheckRepo := healthcheckrepository.NewRepository(db)
	healthcheckService := healthcheckservice.NewService(healthcheckRepo)

	router := routerRestApi{
		handler:                     c,
		TransportHealthCheckRestApi: transporthealthcheck.NewTransportRestApi(healthcheckService),
	}
	router.init()

	return router
}
