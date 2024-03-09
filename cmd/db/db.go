package db

import (
	"fmt"
	"go_pg_pubsub/pkg/models"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/spf13/cobra"
)

func CreatePGMQSchema(db *pg.DB) error {
	query := `CREATE SCHEMA if not exists pgmq`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func CreatePGMQTables(db *pg.DB, models []interface{}) error {
	db.Exec(`set search_path='pgmq'`)

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

func DropPGMQTables(db *pg.DB, models []interface{}) error {
	db.Exec(`set search_path='pgmq'`)

	for _, model := range models {
		err := db.Model(model).DropTable(&orm.DropTableOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

var PGMQModels = []interface{}{
	(*models.PgmqQueue)(nil),
}

// initCmd represents the init command
var InitCmd = &cobra.Command{
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

		err := CreatePGMQSchema(db)
		if err != nil {
			panic(err)
		}

		errTables := CreatePGMQTables(db, PGMQModels)
		if err != nil {
			panic(errTables)
		}
	},
}

// initCmd represents the init command
var CleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running init")
		db := pg.Connect(&pg.Options{
			User:     "postgres",
			Password: "postgres",
			Database: "go_pg_pubsub_dev",
		})
		defer db.Close()

		err := DropPGMQTables(db, PGMQModels)
		if err != nil {
			panic(err)
		}
	},
}
