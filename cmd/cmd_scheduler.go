package main

import (
	"context"
	"erp-directory-service/internal/app"
	"erp-directory-service/internal/config"
	"erp-directory-service/internal/provider"
	"errors"
	"syscall"
	"time"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/confy"
	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/utils/graceful"
	"github.com/spf13/cobra"
)

func newCmdScheduler() *cobra.Command {
	preRunClosed := make([]func() error, 0, 2)
	var shutdownScheduler func(ctx context.Context) error

	var slogHookOption string
	var zerologHookOption string

	cmd := &cobra.Command{
		Use:   "scheduler",
		Short: "Run the scheduler",
		PreRun: func(cmd *cobra.Command, args []string) {
			appCfg := config.GetAppScheduler()

			app.StartPprofServer(cmd.Name())
			closeLogging := provider.NewLogging(
				"scheduler",
				slogHookOption, zerologHookOption,
				appCfg.DebugMode, appCfg.Env, appCfg.Name,
			)

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

				err := shutdownScheduler(ctx)
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
			scheduler := app.NewSchedulerApp()
			shutdownScheduler = scheduler.Shutdown
			scheduler.Start()
		},
	}

	cmd.Flags().StringVarP(&slogHookOption, "slog-hook", "s", "file-writer",
		"slog hook output: file-writer (write to file) or std-out (write to stdout). default: file-writer",
	)
	cmd.Flags().StringVarP(&zerologHookOption, "zerolog-hook", "z", "file-writer",
		"zerolog hook output: file-writer (write to file) or std-out (write to stdout). default: file-writer",
	)

	return cmd
}
