package main

import (
	"github.com/spf13/cobra"
)

func main() {
	root := &cobra.Command{
		Use:   "directoryservice",
		Short: "A CLI Directory Service",
	}

	root.AddCommand(newRestApiCmd())
	root.AddCommand(newGrpcApiCmd())
	err := root.Execute()
	if err != nil {
		panic(err)
	}
}
