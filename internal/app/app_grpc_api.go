package app

import (
	"context"
	"erp-directory-service/internal/config"
	"erp-directory-service/internal/gen/grpcgen/healthcheck"
	"erp-directory-service/internal/infrastructure"
	healthcheckrepository "erp-directory-service/internal/module/healthcheck/repository"
	healthcheckservice "erp-directory-service/internal/module/healthcheck/service"
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcApiApp struct {
	server   *grpc.Server
	port     int
	listener net.Listener
	closeFn  []func() error
}

func NewGrpcApiApp(port int) *grpcApiApp {
	appCfg := config.GetAppGrpcApi()
	if port == 0 {
		port = appCfg.Port
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("gagal listen: %v", err)
	}

	s := grpc.NewServer()

	grpcApp := &grpcApiApp{
		port:     port,
		listener: lis,
		server:   s,
		closeFn:  make([]func() error, 0),
	}

	grpcApp.init()

	return grpcApp
}

func (r *grpcApiApp) Shutdown(ctx context.Context) error {
	errs := make([]error, 0, len(r.closeFn))

	r.server.GracefulStop()
	err := r.listener.Close()
	if err != nil && !errors.Is(err, net.ErrClosed) {
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

func (r *grpcApiApp) Start() {
	slog.Info(fmt.Sprintf("gRPC server running on :%d", r.port))
	if err := r.server.Serve(r.listener); err != nil {
		slog.Error(err.Error())
	}
}

func (r *grpcApiApp) init() {
	db, err := infrastructure.NewDB()
	if err != nil {
		panic(err)
	}
	r.closeFn = append(r.closeFn, db.Close)

	healthcheckRepo := healthcheckrepository.NewRepository(db)
	healthcheckService := healthcheckservice.NewService(healthcheckRepo)

	routerGrpc := routerGrpcApi{
		healthcheck: transporthealthcheck.NewTransportGrpc(healthcheckService),
	}

	routerGrpc.init(r.server)
	reflection.Register(r.server)
}

type routerGrpcApi struct {
	healthcheck *transporthealthcheck.TransportHealthCheckGrpc
}

func (i *routerGrpcApi) init(s *grpc.Server) {
	healthcheck.RegisterHealthCheckServiceServer(s, i.healthcheck)
}
