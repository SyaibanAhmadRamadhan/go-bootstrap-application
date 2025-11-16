package app

import (
	"erp-directory-service/internal/gen/grpcgen/healthcheck"
	transporthealthcheck "erp-directory-service/internal/transport/healthcheck"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

type routerRestApi struct {
	handler *chi.Mux
	// restapigen.Unimplemented
	*transporthealthcheck.TransportHealthCheckRestApi
}

func (r *routerRestApi) init() {
	r.handler.HandleFunc("/debug/pprof/", pprof.Index)
	r.handler.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.handler.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.handler.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.handler.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.handler.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	r.handler.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	r.handler.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	r.handler.Handle("/debug/pprof/block", pprof.Handler("block"))
	r.handler.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	r.handler.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
}

type routerGrpc struct {
	healthcheck *transporthealthcheck.TransportHealthCheckGrpc
}

func (i *routerGrpc) init(s *grpc.Server) {
	healthcheck.RegisterHealthCheckServiceServer(s, i.healthcheck)
}
