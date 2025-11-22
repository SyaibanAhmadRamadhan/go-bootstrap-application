package main

import (
	"context"
	"errors"
	"go-bootstrap/internal/app"
	"go-bootstrap/internal/config"
	"go-bootstrap/internal/infrastructure"
	"syscall"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/graceful"
	"github.com/spf13/cobra"
)

func newRestApiCmd() *cobra.Command {
	preRunClosed := make([]func() error, 0, 1)
	var shutdownServer func(ctx context.Context) error

	var port int
	var slogHookOption string
	var zerologHookOption string

	cmd := &cobra.Command{
		Use:   "restapi",
		Short: "Run the server",
		PreRun: func(cmd *cobra.Command, args []string) {
			app.StartPprofServer()
			closeLogging := infrastructure.NewLogging(slogHookOption, zerologHookOption)

			preRunClosed = append(preRunClosed, closeLogging)
			preRunClosed = append(preRunClosed, config.UnwatchLoader)
			preRunClosed = append(preRunClosed, func() error {
				confy.Close()
				return nil
			})
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
			restApi := app.NewRestApiApp(port)
			shutdownServer = restApi.ShutdownAndClose
			go func() {
				restApi.Start()
			}()
		},
	}

	cmd.Flags().IntVarP(&port, "port", "p", 0, "Port to run the server on")
	cmd.Flags().StringVarP(&slogHookOption, "slog-hook", "s", "file-writer",
		"slog hook output: file-writer (write to file) or std-out (write to stdout). default: file-writer",
	)
	cmd.Flags().StringVarP(&zerologHookOption, "zerolog-hook", "z", "file-writer",
		"zerolog hook output: file-writer (write to file) or std-out (write to stdout). default: file-writer",
	)

	return cmd
}
