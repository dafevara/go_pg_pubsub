/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
    "time"
	"github.com/spf13/cobra"
    "github.com/go-pg/pg/v10"
    "github.com/go-pg/pg/v10/orm"
)

type User struct {
    Id      int32
    Name    string
    Email   string
    Balance int32
}

type Product struct {
    Id          int32
    Name        string
    Price       int32
    Stock       int32
    Discount    int32
}

type Payment struct {
    Id          int32
    ProductId   int32
    Product     *Product `pg:"rel:has-one"`
    UserId      int32
    User        *User `pg:"rel:has-one"`
    Amount      int32
    Status      string
}

type PaymentTask struct {
    Id          int32
    PaymentId   int32
    Payment     *Payment `pg:"rel:has-one"`
    TriesLeft   int32
    Error       string
    Processing  bool
    UpdatedAt   time.Time
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("running init")
        db := pg.Connect(&pg.Options{
            User: "postgres", 
            Password: "postgres",
            Database: "go_pg_pubsub_dev",
        })
        defer db.Close()

        err := createSchema(db)
        if err != nil {
            panic(err)
        }

	},
}

func createSchema(db *pg.DB) error {
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
