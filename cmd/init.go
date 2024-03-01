package cmd

import (
	"fmt"
    "go_pg_pubsub/pkg/types"

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

		err = AddTriggers(db)
		if err != nil {
			panic(err)
		}
	},
}

func AddTriggers(db *pg.DB) error {
	_, err := db.Exec(`
        CREATE OR REPLACE FUNCTION insert_into_payment_task()
        RETURNS TRIGGER AS $$
        BEGIN
        INSERT INTO payment_tasks (payment_id) VALUES (NEW.id);
        RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;

        drop trigger if exists process_payment_trigger on payments cascade;
        CREATE TRIGGER process_payment_trigger
        AFTER INSERT
        ON payments
        FOR EACH ROW
        EXECUTE FUNCTION insert_into_payment_task();
    `)
	if err != nil {
		return err
	}

	return nil
}

func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*types.User)(nil),
		(*types.Product)(nil),
		(*types.Payment)(nil),
		(*types.PaymentTask)(nil),
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
