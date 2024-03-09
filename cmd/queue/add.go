package queue

import (
	"github.com/go-pg/pg/v10"
	"github.com/spf13/cobra"
)

func Add(queueName string) error {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()

	return nil
}

func List() error {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()

	return nil
}

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "",
	Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		err := Add(args[0])
		if err != nil {
			panic(err)
		}
	},
}

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "",
	Run: func(cmd *cobra.Command, args []string) {
		err := List()
		if err != nil {
			panic(err)
		}
	},
}
