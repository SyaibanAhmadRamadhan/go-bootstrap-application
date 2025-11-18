package main

import (
	"erp-directory-service/internal/config"
	"log"

	"github.com/spf13/cobra"
)

func init() {

}

func main() {
	root := &cobra.Command{
		Use:   "directoryservice",
		Short: "A CLI Directory Service",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			switch cmd.Name() {
			case "scheduler":
				config.LoadConfig("app_scheduler.debug_mode")
			case "restapi":
				config.LoadConfig("app_rest_api.debug_mode")
			case "grpcapi":
				config.LoadConfig("app_grpc_api.debug_mode")
			}
		},
	}

	root.AddCommand(newRestApiCmd())
	root.AddCommand(newGrpcApiCmd())
	root.AddCommand(newCmdScheduler())

	err := root.Execute()
	if err != nil {
		log.Fatal(err.Error())
	}
}
