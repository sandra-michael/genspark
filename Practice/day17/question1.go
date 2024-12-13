package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	_ "github.com/jackc/pgx/v5/stdlib"

)

// q1. Create docker compose file to run  postgres container
//     Connect to postgres using pgx
//     Create movies table using Go program
//     Insert two records within transaction to movies table
//     Update one record using optimistic locking

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

var DB *sql.DB

func main() {
	defer DB.Close()
	err := DB.Ping()
	if err != nil {
		panic(err)
	}

	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()
	err = Create(tx)
	if err != nil {
		log.Println(err)
		return
	}

	err = InsertQuery(tx)
	if err != nil {
		log.Println(err)
		return
	}

	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
}

func Create(tx *sql.Tx) error {
	createQuery := `CREATE TABLE IF NOT EXISTS movies (
	id SERIAL PRIMARY KEY,
	name text NOT NULL,
	updated_at TIMESTAMP NOT NULL
	);`

	_, err := tx.Exec(createQuery)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func InsertQuery(tx *sql.Tx) error {
	// Prepare the insert queries
	insertQuery1 := `INSERT INTO movies (name, updated_at) VALUES ($1, $2);`
	insertQuery2 := `INSERT INTO movies (name, updated_at) VALUES ($1, $2);`

	// Data for first author
	name1 := "Men In Black"

	// Data for second author
	name2 := "Die Hard"

	// Execute the first insert query
	_, err := tx.Exec(insertQuery1, name1, time.Now().UTC())
	if err != nil {
		log.Printf("Failed to insert first movie: %v", err)
		return err
	}

	// Execute the second insert query
	_, err = tx.Exec(insertQuery2, name2, time.Now().UTC())
	if err != nil {
		log.Printf("Failed to insert second movie: %v", err)
		return err
	}
	return nil
}

func UpdateQuery(tx *sql.Tx, id int) error {
	var updatedAt time.Time

	err := tx.QueryRow(`
		UPDATE movies
		SET name = $1, updated_at = $2
		WHERE id = $3$
		RETURNING updated_at;`,
		"Men In BLack 2", time.Now().UTC(), id).Scan(&updatedAt)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
