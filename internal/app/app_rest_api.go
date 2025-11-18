package app

import (
	"context"
	"erp-directory-service/internal/config"
	"erp-directory-service/internal/gen/restapigen"
	"erp-directory-service/internal/infrastructure"
	healthcheckrepository "erp-directory-service/internal/module/healthcheck/repository"
	healthcheckservice "erp-directory-service/internal/module/healthcheck/service"
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/chix"
	"github.com/go-chi/chi/v5"
)

type restApiApp struct {
	server  *http.Server
	port    int
	closeFn []func() error
}

func NewRestApiApp(port int) *restApiApp {
	appCfg := config.GetAppRestApi()
	if port == 0 {
		port = appCfg.Port
	}

	handler := chix.New(chix.Config{
		BlacklistRouteLogResponse: map[string]struct{}{},
		SensitiveFields:           map[string]struct{}{},
		CorsConf: chix.CorsConfig{
			AllowOrigins:     nil,
			AllowMethods:     nil,
			AllowHeaders:     nil,
			AllowCredentials: true,
		},
		AppName: appCfg.Name,
		UseOtel: false,
	})

	restapiApp := &restApiApp{
		port:    port,
		closeFn: make([]func() error, 0),
	}

	restapiApp.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: restapigen.HandlerFromMux(restapiApp.init(handler), handler),
	}

	return restapiApp
}

func (r *restApiApp) ShutdownAndClose(ctx context.Context) error {
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

func (r *restApiApp) Start() {
	slog.Info(fmt.Sprintf("REST API listening on %s", r.server.Addr))
	err := r.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error(err.Error())
	}
}

func (r *restApiApp) init(c *chi.Mux) routerRestApi {
	db, err := infrastructure.NewDB()
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

	return router
}

type routerRestApi struct {
	handler *chi.Mux
	// restapigen.Unimplemented
	*transporthealthcheck.TransportHealthCheckRestApi
}
