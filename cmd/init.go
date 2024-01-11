package cmd

import (
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running init")
		db := pg.Connect(&pg.Options{
			User:     "postgres",
			Password: "postgres",
			Database: "go_pg_pubsub_dev",
		})
		defer db.Close()

		err := CreateSchema(db)
		if err != nil {
			panic(err)
		}
	},
}

func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Product)(nil),
		(*Payment)(nil),
		(*PaymentTask)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)
}
