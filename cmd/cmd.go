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
			config.LoadConfig(cmd.Name())
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
