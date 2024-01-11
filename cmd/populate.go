/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/bxcodec/faker/v4"
	"github.com/go-pg/pg/v10"
	"github.com/spf13/cobra"
)

type FakeUser struct {
	Name    string `faker:"name"`
	Email   string `faker:"email"`
	Balance int64  `faker:"boundary_start=1000, boundary_end=100000"`
}

type FakeProduct struct {
	Name     string `faker:"word"`
	Price    int64  `faker:"boundary_start=10, boundary_end=10000"`
	Stock    int32  `faker:"boundary_start=0, boundary_end=100"`
	Discount int32  `faker:"boundary_start=0, boundary_end=50"`
}

func Populate(db *pg.DB) error {
	for i := 0; i < 100; i++ {
		fUser := FakeUser{}
		err := faker.FakeData(&fUser)
		if err != nil {
			fmt.Println(err)
		}

		user := &User{
			Name:    fUser.Name,
			Email:   fUser.Email,
			Balance: fUser.Balance,
		}

		_, err = db.Model(user).Insert()
		if err != nil {
			panic(err)
		}

		fProduct := FakeProduct{}
		err := faker.FakeData(&fProduct)
		if err != nil {
			fmt.Println(err)
		}

		product := &Product{
			Name:     fProduct.Name,
			Price:    fProduct.Price,
			Stock:    fProduct.Stock,
			Discount: fProduct.Discount,
		}

		_, err = db.Model(product).Insert()
		if err != nil {
			panic(err)
		}
	}
	return nil
}

// populateCmd represents the populate command
var populateCmd = &cobra.Command{
	Use:   "populate",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("populate called")
		db := pg.Connect(&pg.Options{
			User:     "postgres",
			Password: "postgres",
			Database: "go_pg_pubsub_dev",
		})
		defer db.Close()

		err := Populate(db)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(populateCmd)
}
