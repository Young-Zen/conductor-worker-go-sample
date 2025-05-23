/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"worker-sample/conductor"
	"worker-sample/config"
	"worker-sample/database/mysql"
	"worker-sample/server"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Conductor worker sample application",
	Long:  `Run Conductor worker sample application.`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := config.NewServiceContext(configPath, configName)

		mysql.InitDB(ctx)
		defer mysql.CloseDB(ctx)

		conductor.InitWorker(ctx)

		server.InitHttpServer(ctx)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
