/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tasks

import (
	"fmt"
	"go_pg_pubsub/pkg/models"

	"github.com/go-faker/faker/v4"
	"github.com/go-pg/pg/v10"
	"github.com/schollz/progressbar/v3"
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

		payment := &models.Payment{
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
