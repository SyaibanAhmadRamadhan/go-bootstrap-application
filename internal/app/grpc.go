package app

import (
	"context"
	healthcheckrepository "erp-directory-service/internal/module/healthcheck/repository"
	healthcheckservice "erp-directory-service/internal/module/healthcheck/service"
	"erp-directory-service/internal/provider"
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcApi struct {
	server   *grpc.Server
	port     int
	listener net.Listener
	closeFn  []func() error
}

func NewGrpcApi(port int) *grpcApi {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("gagal listen: %v", err)
	}

	s := grpc.NewServer()

	grpcApi := &grpcApi{
		port:     port,
		listener: lis,
		server:   s,
		closeFn:  make([]func() error, 0),
	}

	grpcApi.init()

	return grpcApi
}

func (r *grpcApi) Shutdown(ctx context.Context) error {
	errs := make([]error, 0, len(r.closeFn))

	r.server.GracefulStop()
	err := r.listener.Close()
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

func (r *grpcApi) Start() {
	slog.Info(fmt.Sprintf("gRPC server running on :%d", r.port))
	if err := r.server.Serve(r.listener); err != nil {
		slog.Error(err.Error())
	}
}

func (r *grpcApi) init() {
	db, err := provider.NewDB()
	if err != nil {
		panic(err)
	}
	r.closeFn = append(r.closeFn, db.Close)

	healthcheckRepo := healthcheckrepository.NewRepository(db)
	healthcheckService := healthcheckservice.NewService(healthcheckRepo)

	routerGrpc := routerGrpc{
		healthcheck: transporthealthcheck.NewTransportGrpc(healthcheckService),
	}

	routerGrpc.init(r.server)
	reflection.Register(r.server)
}
