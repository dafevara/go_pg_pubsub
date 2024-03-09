package models

import (
	"go_pg_pubsub/pkg/interfaces"
)

type User struct {
	Id      int32
	Name    string
	Email   string
	Balance int64
}

type Product struct {
	Id       int32
	Name     string
	Price    int64
	Stock    int32
	Discount int32
}

type Payment struct {
	Id        int32
	ProductId int32
	Product   *Product `pg:"rel:has-one"`
	UserId    int32
	User      *User `pg:"rel:has-one"`
	Amount    int64
	Status    string
}

type PaymentTask struct {
	interfaces.BaseTaskable
	PaymentId int32
	Payment   *Payment `pg:"rel:has-one"`
}

type PgmqQueue struct {
	Id   int32
	Name string
}
