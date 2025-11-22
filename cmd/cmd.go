package main

import (
	"go-bootstrap/internal/config"
	"log"

	"github.com/SyaibanAhmadRamadhan/go-foundation-kit/validatorx"
	"github.com/spf13/cobra"
)

func init() {

}

func main() {
	root := &cobra.Command{
		Use:   "go-boostrap",
		Short: "A CLI Golang Boostrap",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			config.LoadConfig(cmd.Name())
			validatorx.InitValidator()
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
