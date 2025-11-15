package main

import (
	"context"
	"erp-directory-service/internal/app"
	"erp-directory-service/internal/config"
	"erp-directory-service/internal/provider"
	"errors"
	"syscall"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/graceful"
	"github.com/spf13/cobra"
)

var port int
var slogHook string
var zerologHook string

func newRestApiCmd() *cobra.Command {
	preRunClosed := make([]func() error, 0, 1)
	var shutdownServer func(ctx context.Context) error

	cmd := &cobra.Command{
		Use:   "restapi",
		Short: "Run the server",
		PreRun: func(cmd *cobra.Command, args []string) {
			preRunClosed = append(preRunClosed, provider.NewLogging(slogHook, zerologHook))
		},
		PostRun: func(cmd *cobra.Command, args []string) {
			graceful.Shutdown(func(ctx context.Context) error {
				errs := make([]error, 0)

				err := shutdownServer(ctx)
				if err != nil {
					errs = append(errs, err)
				}

				for _, v := range preRunClosed {
					if err := v(); err != nil {
						errs = append(errs, err)
					}
				}

				return errors.Join(errs...)
			}, 30*time.Second, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		},
		Run: func(cmd *cobra.Command, args []string) {
			restApi := app.NewRestApi(port)
			shutdownServer = restApi.Shutdown
			go func() {
				restApi.Start()
			}()
		},
	}

	config.LoadConfig()
	cmd.Flags().IntVarP(&port, "port", "p", config.GetApp().Port, "Port to run the server on")
	cmd.Flags().StringVarP(&slogHook, "slog-hook", "s", "file-writer",
		"slog hook output: file-writer (write to file) or std-out (write to stdout). default: file-writer",
	)
	cmd.Flags().StringVarP(&zerologHook, "zerolog-hook", "z", "file-writer",
		"zerolog hook output: file-writer (write to file) or std-out (write to stdout). default: file-writer",
	)

	return cmd
}
