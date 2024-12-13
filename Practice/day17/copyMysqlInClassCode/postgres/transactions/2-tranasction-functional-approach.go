package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func init() {
	const (
		host     = "localhost"
		port     = "5433"
		user     = "postgres"
		password = "postgres"
		dbname   = "postgres"
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	DB, err = sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

}

// insert two more records within transaction

func main() {
	err := DB.Ping()
	if err != nil {
		panic(err)
	}

	err = Update()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Done")
}

func Update() error {

	// wrapping our queries inside the func that takes tx as param
	// we can pass this func to withTx
	f := func(tx *sql.Tx) error {
		updateQuery := `UPDATE author
					SET name = $1
					WHERE email = $2;`

		_, err := tx.Exec(updateQuery, "ABC", "john.doe@example.com")
		if err != nil {
			return fmt.Errorf("update error: %w", err)
		}

		// _, err = tx.Exec(updateQuery, "John1", "john.doe@example.com")
		// if err != nil {
		// 	log.Println(err)
		// 	return fmt.Errorf("update error: %w", err)
		// }
		return nil
	}

	err := withTx(context.Background(), f)
	if err != nil {
		return err
	}
	return nil
}

// withTx func takes a context, and a function that want to exec within a transaction
func withTx(ctx context.Context, fn func(*sql.Tx) error) error {

	// begin a transaction
	tx, err := DB.BeginTx(ctx, nil)

	// call the function passed to withTX,
	// passing the newly created transaction
	err = fn(tx) // func would use tx to run queries within transaction
	if err != nil {
		// rollback in case if any error happens
		er := tx.Rollback()
		if er != nil {
			return fmt.Errorf("rollback error: %w", err)
		}
		return fmt.Errorf("transaction error: %w", err)
	}

	// commit if no error
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit error: %w", err)
	}
	return nil

}
