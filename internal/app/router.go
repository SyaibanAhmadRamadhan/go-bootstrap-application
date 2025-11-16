package app

import (
	"erp-directory-service/gen/grpcgen/healthcheck"
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"

	"google.golang.org/grpc"
)

type routerRestApi struct {
	// restapigen.Unimplemented
	*transporthealthcheck.TransportHealthCheckRestApi
}

type routerGrpc struct {
	healthcheck *transporthealthcheck.TransportHealthCheckGrpc
}

func (i *routerGrpc) init(s *grpc.Server) {
	healthcheck.RegisterHealthCheckServiceServer(s, i.healthcheck)
}
