/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package queue

import (
	"fmt"
	"go_pg_pubsub/pkg/models"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/spf13/cobra"
)

func Next() (models.PaymentTask, error) {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()
	query := `
        UPDATE payment_tasks SET
            processing = true,
            tries_left = tries_left - 1,
            error = NULL,
            next_try_at = NULL,
            updated_at = CURRENT_TIMESTAMP
        WHERE id = (
            SELECT id
            FROM payment_tasks
            WHERE tries_left > 0
            AND (
                next_try_at IS NULL OR
                next_try_at < CURRENT_TIMESTAMP
            )
            AND (
                processing = false OR
                updated_at < CURRENT_TIMESTAMP - INTERVAL '1 SEC'
            )
            ORDER BY next_try_at ASC, id ASC
            FOR UPDATE SKIP LOCKED
            LIMIT 1
        )
        RETURNING id, payment_id, tries_left, error, processing, updated_at
    `
	var paymentTask models.PaymentTask
	_, err := db.QueryOne(&paymentTask, query)
	if err != nil {
		return paymentTask, err
	}

	return paymentTask, nil
}

func Success(paymentTask *models.PaymentTask, payment models.Payment, product models.Product, newBalance int64) error {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()
	// Update user's balance
	q := "update users set balance = ? where id = ?"
	_, err := db.Exec(q, newBalance, payment.UserId)
	if err != nil {
		fmt.Println("Error updating user balance: ", err)
		return err
	}

	// Update payment with success state
	q = "update payments set status = ? where id = ?"
	_, err = db.Exec(q, "accepted", paymentTask.PaymentId)
	if err != nil {
		fmt.Println("Error updating payment status: ", err)
		return err
	}

	fmt.Println("Successful payment. New Balance: ", newBalance)

	// Update product with new stock
	newStock := product.Stock - 1
	fmt.Printf("New Stock: %v\n", newStock)
	q = "update products set stock = ? where id = ?"
	_, err = db.Exec(q, newStock, product.Id)
	if err != nil {
		fmt.Println("Error updating product stock: ", err)
		return err
	}

	// Remove message from queue
	query := "delete from payment_tasks where id = ?"
	_, err = db.Exec(query, paymentTask.Uuid)
	if err != nil {
		fmt.Println("Error deleting payment task: ", err)
		return err
	}

	return nil
}

func FailedByStock(paymentTask *models.PaymentTask, user models.User, payment models.Payment) error {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()

	msg := "Unable to pay because there's no stock"
	fmt.Printf("Failed: %v\n", msg)

	q := "update payment_tasks set error = ?, tries_left = 0, processing = false where id = ?"
	_, err := db.Exec(q, msg, paymentTask.Uuid)
	if err != nil {
		fmt.Printf("Error updating payment task: %v", err)
		return err
	}

	q = "update payments set status = 'rejected' where id = ?"
	_, err = db.Exec(q, payment.Id)
	if err != nil {
		fmt.Printf("Error updating payment status: %v", err)
		return err
	}

	return nil
}

func FailedByBalance(paymentTask *models.PaymentTask, user models.User, payment models.Payment, product models.Product) error {
	balance := user.Balance
	price := product.Price
	msg := fmt.Sprintf("Unable to pay because price: %d is greater than balance %d", price, balance)

	fmt.Printf("Failed: %v\n", msg)

	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()

	q := "update payment_tasks set error = ? where id = ?"
	_, err := db.Exec(q, msg, paymentTask.Uuid)
	if err != nil {
		fmt.Printf("Error updating payment task: %v", err)
		return err
	}

	q = "update payments set status = 'rejected' where id = ?"
	_, err = db.Exec(q, payment.Id)
	if err != nil {
		fmt.Printf("Error updating payment status: %v", err)
		return err
	}

	return nil
}

func Perform(paymentTask *models.PaymentTask) error {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()
	// Assuming db.GetClient returns a *sql.DB and handles error internally
	// Query Payment
	var payment models.Payment
	q := "select * from payments where id = ?"
	_, err := db.QueryOne(&payment, q, paymentTask.PaymentId)
	if err != nil {
		return err
	}
	fmt.Printf("Payment: %v\n", payment.Id)

	// Query User
	var user models.User
	q = "select * from users where id = ?"
	_, err = db.QueryOne(&user, q, payment.UserId)
	if err != nil {
		return err
	}
	fmt.Printf("for User: %v\n", user.Id)

	// Query Product
	var product models.Product
	q = "select * from products where id = ?"
	_, err = db.QueryOne(&product, q, payment.ProductId)
	if err != nil {
		return err
	}
	fmt.Printf("Product %v\n", product.Id)

	balance := user.Balance
	price := product.Price
	stock := product.Stock

	newBalance := balance - price
	newStock := stock - 1
	fmt.Println("Balance: ", balance)
	fmt.Println("Price: ", price)
	fmt.Println("Stock: ", stock)
	fmt.Println("New Balance: ", newBalance)
	fmt.Println("New Stock: ", newStock)

	if newBalance < 0 {
		return FailedByBalance(paymentTask, user, payment, product) // assuming this function is defined
	} else if newStock < 0 {
		return FailedByStock(paymentTask, user, payment) // assuming this function is defined
	} else {
		return Success(paymentTask, payment, product, newBalance) // assuming this function is defined
	}

	return nil
}

func Attach() error {
	db := pg.Connect(&pg.Options{
		User:     "postgres",
		Password: "postgres",
		Database: "go_pg_pubsub_dev",
	})
	defer db.Close()
	for {
		if next, err := Next(); err == nil {
			fmt.Println(next)
			err = Perform(&next)
			if err != nil {
				fmt.Println(err)
			}
			// time.Sleep(1 * time.Second)
		} else {
			fmt.Println("Error:", err)
			time.Sleep(5 * time.Second)
		}
	}
}

// subscribeCmd represents the subscribe command
var subscribeCmd = &cobra.Command{
	Use:   "subscribe",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("subscribe called")
		err := Attach()
		if err != nil {
			panic(err)
		}
	},
}
