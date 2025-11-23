package app

import (
	"context"
	"errors"
	"fmt"
	"go-bootstrap/internal/config"
	"go-bootstrap/internal/gen/restapigen"
	"go-bootstrap/internal/infrastructure"
	authrepository "go-bootstrap/internal/module/auth/repository"
	authservice "go-bootstrap/internal/module/auth/service"
	healthcheckrepository "go-bootstrap/internal/module/healthcheck/repository"
	healthcheckservice "go-bootstrap/internal/module/healthcheck/service"
	userrepository "go-bootstrap/internal/module/user/repository"
	userservice "go-bootstrap/internal/module/user/service"
	transportauth "go-bootstrap/internal/transport/auth"
	transporthealthcheck "go-bootstrap/internal/transport/healthcheck"
	transportuser "go-bootstrap/internal/transport/user"
	"io"
	"log/slog"
	"net/http"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/http/server/ginx"
	"github.com/gin-gonic/gin"
)

type restApiApp struct {
	server    *http.Server
	ginEngine *gin.Engine
	port      int
	closeFn   []func() error
}

func NewRestApiApp(port int) *restApiApp {
	appCfg := config.GetAppRestApi()
	if port == 0 {
		port = appCfg.Port
	}

	if appCfg.Gin.DisableDefaultWriter {
		gin.DefaultWriter = io.Discard
	}
	if appCfg.Gin.DisableErrorWriter {
		gin.DefaultErrorWriter = io.Discard
	}

	ginEngine := ginx.NewGin(ginx.GinConfig{
		BlacklistRouteLogResponse: map[string]struct{}{},
		SensitiveFields:           map[string]struct{}{},
		CorsConf: ginx.CorsConfig{
			AllowOrigins:     appCfg.Gin.Cors.AllowOrigins,
			AllowMethods:     appCfg.Gin.Cors.AllowMethods,
			AllowHeaders:     appCfg.Gin.Cors.AllowHeaders,
			AllowCredentials: appCfg.Gin.Cors.AllowCredentials,
			ExposeHeaders:    appCfg.Gin.Cors.ExposeHeaders,
			MaxAge:           appCfg.Gin.Cors.MaxAge,
		},
		AppName: appCfg.Name,
		UseOtel: appCfg.Gin.UseOtel,
	})

	gin.SetMode(appCfg.Gin.Mode)
	if appCfg.Gin.DisableConsoleColor {
		gin.DisableConsoleColor()
	}

	restapiApp := &restApiApp{
		port:      port,
		ginEngine: ginEngine,
		closeFn:   make([]func() error, 0),
	}

	restapigen.RegisterHandlers(ginEngine, restapiApp.init())

	restapiApp.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: ginEngine,
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

func (r *restApiApp) init() routerRestApi {
	db, err := infrastructure.NewDB()
	if err != nil {
		panic(err)
	}
	r.closeFn = append(r.closeFn, db.Close)

	ginHelper := ginx.NewGinHelper("message", "errors")

	healthcheckService := healthcheckservice.NewService(
		healthcheckrepository.NewRepository(db),
	)

	authService := authservice.NewService(
		authrepository.NewRepository(db),
		authrepository.NewUserRepository(db),
	)

	userService := userservice.NewService(
		userrepository.NewRepository(db),
	)

	router := routerRestApi{
		HealthCheckRestApiHandler: transporthealthcheck.NewRestApiHandler(healthcheckService),
		AuthRestAPIHandler:        transportauth.NewRestAPIHandler(authService, ginHelper),
		UserRestAPIHandler:        transportuser.NewRestAPIHandler(userService, ginHelper),
	}

	return router
}

type routerRestApi struct {
	// restapigen.Unimplemented
	*transporthealthcheck.HealthCheckRestApiHandler
	*transportauth.AuthRestAPIHandler
	*transportuser.UserRestAPIHandler
}
