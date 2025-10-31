package main

import (
	"erp-directory-service/internal/config"
	"erp-directory-service/internal/provider"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var port int

func init() {
	config.LoadConfig()
	provider.NewLogging()
	restApiCmd.Flags().IntVarP(&port, "port", "p", config.GetApp().Port, "Port to run the server on")
}

var restApiCmd = &cobra.Command{
	Use:   "restapi",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		time.Sleep(5 * time.Second)
		fmt.Println("Server running...")
	},
}
