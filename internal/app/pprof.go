package app

import (
	"context"
	"erp-directory-service/internal/config"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"sync/atomic"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy"
)

func StartPprofServer(cmdName string) {
	handler := http.NewServeMux()
	registerPprof(handler)

	httpServer := &http.Server{
		Handler: handler,
	}
	isRunning := uint32(0)

	pprofConfig := getPprofConfig(cmdName)
	if pprofConfig.Enable && atomic.LoadUint32(&isRunning) == 0 {
		httpServer.Addr = fmt.Sprintf(":%d", pprofConfig.Port)

		go func() {
			slog.Info("bootstrap application, pprof listening", "addr", httpServer.Addr)
			if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				slog.Error("bootstrap applicaiton, pprof listen error", "err", err)
			}
		}()

		atomic.StoreUint32(&isRunning, 1)
	}

	_, confySubsribetionSignal := confy.Subscribe()
	go func() {
		for range confySubsribetionSignal {
			pprofConfig := getPprofConfig(cmdName)
			if pprofConfig.Enable && atomic.LoadUint32(&isRunning) == 0 {
				httpServer.Addr = fmt.Sprintf(":%d", pprofConfig.Port)

				go func() {
					slog.Info("pprof listening", "addr", httpServer.Addr)
					if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						slog.Error("pprof listen error", "err", err)
					}
				}()
				atomic.StoreUint32(&isRunning, 1)
			} else if !pprofConfig.Enable && atomic.LoadUint32(&isRunning) == 1 {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				err := httpServer.Shutdown(ctx)
				if err != nil {
					slog.Error("pprof shutdown server error", "err", err)
				}
				cancel()
				atomic.StoreUint32(&isRunning, 0)
				slog.Info("pprof shutdown server successfully")
			}
		}
	}()

}

func registerPprof(mux *http.ServeMux) {
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))
	mux.Handle("/debug/pprof/allocs", pprof.Handler("allocs"))
	mux.Handle("/debug/pprof/mutex", pprof.Handler("mutex"))
}

func getPprofConfig(cmdName string) config.Pprof {
	switch cmdName {
	case "scheduler":
		return config.GetPprofAppScheduler()
	case "restapi":
		return config.GetPprofAppRestApi()
	case "grpcapi":
		return config.GetPprofAppGrpcApi()
	default:
		slog.Error("unknown cmd name for get pprof config")
		return config.Pprof{}
	}
}
