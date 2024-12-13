package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	"time"
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

func main() {

	tx, err := DB.Begin()
	if err != nil {
		log.Println(err)
		return
	}
	defer tx.Rollback()

	var id, version int
	err = tx.QueryRow(`SELECT id, version FROM students WHERE id = $1`, 1).
		Scan(&id, &version)

	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("Existing Record: ID=%d,, Version=%d\n",
		id, version)

	time.Sleep(10 * time.Second)
	newName := "New Name"
	var updatedAt time.Time
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
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("Record Updated: ID=%d, New Name=%s, UpdatedAt=%s\n",
		id, newName, updatedAt.Format(time.RFC850))

}
