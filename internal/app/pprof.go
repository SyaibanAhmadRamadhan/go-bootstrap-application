package app

import (
	"context"
	"errors"
	"fmt"
	"go-bootstrap/internal/config"
	"log/slog"
	"net/http"
	"net/http/pprof"
	"sync/atomic"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy"
)

func StartPprofServer() {
	handler := http.NewServeMux()
	registerPprof(handler)

	httpServer := &http.Server{
		Handler: handler,
	}
	isRunning := uint32(0)

	pprofConfig := config.GetPprof()
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
			pprofConfig := config.GetPprof()
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
	mux.Handle("/debug/pprof/", authMiddlewarePprof(http.HandlerFunc(pprof.Index)))
	mux.Handle("/debug/pprof/cmdline", authMiddlewarePprof(http.HandlerFunc(pprof.Cmdline)))
	mux.Handle("/debug/pprof/profile", authMiddlewarePprof(http.HandlerFunc(pprof.Profile)))
	mux.Handle("/debug/pprof/symbol", authMiddlewarePprof(http.HandlerFunc(pprof.Symbol)))
	mux.Handle("/debug/pprof/trace", authMiddlewarePprof(http.HandlerFunc(pprof.Trace)))

	mux.Handle("/debug/pprof/heap", authMiddlewarePprof(pprof.Handler("heap")))
	mux.Handle("/debug/pprof/goroutine", authMiddlewarePprof(pprof.Handler("goroutine")))
	mux.Handle("/debug/pprof/threadcreate", authMiddlewarePprof(pprof.Handler("threadcreate")))
	mux.Handle("/debug/pprof/block", authMiddlewarePprof(pprof.Handler("block")))
	mux.Handle("/debug/pprof/allocs", authMiddlewarePprof(pprof.Handler("allocs")))
	mux.Handle("/debug/pprof/mutex", authMiddlewarePprof(pprof.Handler("mutex")))
}

func authMiddlewarePprof(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authotization")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		if token != config.GetPprof().StaticToken {
			http.Error(w, "invalid token", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
