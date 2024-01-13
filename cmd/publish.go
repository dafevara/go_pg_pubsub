/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/go-faker/faker/v4"
	"github.com/go-pg/pg/v10"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

type FakePayment struct {
	ProductId int32 `faker:"boundary_start=1, boundary_end=100"`
	UserId    int32 `faker:"boundary_start=1, boundary_end=100"`
	Amount    int64 `faker:"boundary_start=10, boundary_end=1000"`
}

func Publish(db *pg.DB) error {
	bar := progressbar.Default(1000)
	for i := 0; i < 1000; i++ {
		fPayment := FakePayment{}
		err := faker.FakeData(&fPayment)
		if err != nil {
			fmt.Println(err)
		}

		payment := &Payment{
			ProductId: fPayment.ProductId,
			UserId:    fPayment.UserId,
			Amount:    fPayment.Amount,
		}

		_, err = db.Model(payment).Insert()
		if err != nil {
			panic(err)
		}
		bar.Add(1)
	}
	return nil
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

		err := Publish(db)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// publishCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// publishCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
