package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// https://stackoverflow.com/questions/129329/optimistic-vs-pessimistic-locking/129397#129397

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

func main() {

	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()

	createTable := `CREATE TABLE IF NOT EXISTS students (
	    id SERIAL PRIMARY KEY,
	    name TEXT NOT NULL,
	    email TEXT UNIQUE NOT NULL,
	    updated_at TIMESTAMP NOT NULL,
	    version integer NOT NULL DEFAULT 1
	);`
	_, err = tx.Exec(createTable)

	if err != nil {
		log.Println(err)
		return
	}

	var id, version int
	// grabbing the version number as well with the id
	err = tx.QueryRow(`SELECT id, version FROM students WHERE id = $1`, 1).
		Scan(&id, &version)

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Existing Record: ID=%d,, Version=%d\n",
		id, version)

	newName := "ABC"
	var updatedAt time.Time

	// in this query we check if version number is changed from last select then
	// this update would not work due to & condition in the where clause
	err = tx.QueryRow(`
		UPDATE students
		SET name = $1, updated_at = $2, version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version, updated_at;`,
		newName, time.Now().UTC(), id, version).Scan(&version, &updatedAt)
	if err != nil {
		log.Println(err)
		return
	}

	// if no problem, we will commit the transaction
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Record Updated: ID=%d, New Name=%s, UpdatedAt=%s\n",
		id, newName, updatedAt.Format(time.RFC850))

}
