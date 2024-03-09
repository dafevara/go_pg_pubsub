/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"go_pg_pubsub/cmd/db"
	"go_pg_pubsub/cmd/queue"
	"go_pg_pubsub/cmd/tasks"
	"os"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_pg_pubsub",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running publish")
		db := pg.Connect(&pg.Options{
			User:     "postgres",
			Password: "postgres",
			Database: "go_pg_pubsub_dev",
		})
		defer db.Close()

		err := tasks.Publish(db)
		if err != nil {
			panic(err)
		}
	},
}

// queueCmd represents the queue command
var queueCmd = &cobra.Command{
	Use:   "queue",
	Short: "A brief description of your command",
}

// dbCmd represents the db command
var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.AddCommand(publishCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(dbCmd)

	dbCmd.AddCommand(db.InitCmd)
	dbCmd.AddCommand(db.CleanCmd)

	queueCmd.AddCommand(queue.AddCmd)
	queueCmd.AddCommand(queue.ListCmd)
}
