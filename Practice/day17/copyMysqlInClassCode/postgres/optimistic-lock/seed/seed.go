package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var DB *sql.DB

func main() {
	// Connect to the database
	defer DB.Close()

	createTable := `CREATE TABLE IF NOT EXISTS students (
	    id SERIAL PRIMARY KEY,
	    name TEXT NOT NULL,
	    email TEXT UNIQUE NOT NULL,
	    updated_at TIMESTAMP NOT NULL,
	    version integer NOT NULL DEFAULT 1
	);`
	_, err := DB.Exec(createTable)

	if err != nil {
		log.Println(err)
		return
	}
	// Prepare the insert queries
	insertQuery1 := `INSERT INTO students (name, email, updated_at) VALUES ($1, $2, $3);`
	insertQuery2 := `INSERT INTO students (name, email, updated_at) VALUES ($1, $2, $3);`

	// Data for first author
	name1 := "John Doe"
	email1 := "johndoe@example.com"

	// Data for second author
	name2 := "Jane Smith"
	email2 := "janesmith@example.com"

	// Execute the first insert query
	_, err = DB.Exec(insertQuery1, name1, email1, time.Now().UTC())
	if err != nil {
		log.Printf("Failed to insert first students: %v", err)
		return
	}

	// Execute the second insert query
	_, err = DB.Exec(insertQuery2, name2, email2, time.Now().UTC())
	if err != nil {
		log.Printf("Failed to insert second students: %v", err)
		return
	}

	log.Println("Both students inserted successfully.")
}

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
