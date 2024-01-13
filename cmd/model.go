package cmd

import (
	"time"
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
	Id         int32
	PaymentId  int32
	Payment    *Payment `pg:"rel:has-one"`
	TriesLeft  int32    `pg:"default:5"`
	Error      string
	Processing bool `pg:"default:false"`
	UpdatedAt  time.Time
}
