package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DB as global var not recommended
var DB *sql.DB

// this is used to initialize the state for the current package
// not recommend to be used most of the times.
// hard to test, hard to know when it runs
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
		panic(err)
	}
}

func main() {
	UpdateAuthor()
	err := DB.Ping()
	if err != nil {
		// if db is not connected, no point to continue
		panic(err)
	}

}

func UpdateAuthor() {
	//BeginTx would start the transaction
	tx, err := DB.BeginTx(context.Background(), nil)
	if err != nil {
		log.Println(err)
		return
	}

	// calling rollback multiple times have no effect after commit
	// rollback would roll back any changes if function return early without commit
	defer func() {
		err := tx.Rollback()
		if err != nil {
			log.Println(err)
			return
		}
	}()

	// createQuery := `CREATE TABLE IF NOT EXISTS author (
	// id SERIAL PRIMARY KEY,
	// name text NOT NULL,
	// email text UNIQUE NOT NULL

	// );`

	// _, err = tx.Exec(createQuery)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	// 	insertQuery1 := `INSERT INTO author (name, email)
	// VALUES ('John Doe', 'john.doe@example.com');`

	// 	_, err = tx.Exec(insertQuery1)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	insertQuery2 := `INSERT INTO author (name, email)
	// VALUES ('Jane Smith', 'jane.smith@example.com');`

	// 	_, err = tx.Exec(insertQuery2)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	updateQuery := `UPDATE author
					SET name = $1
					WHERE email = $2;`

	_, err = tx.Exec(updateQuery, "ABC", "john.doe@example.com")
	if err != nil {
		log.Println(err)
		return
	}

	_, err = tx.Exec(updateQuery, "John1", "john.doe@example.com")
	if err != nil {
		log.Println(err)
		return
	}

	// only if both transaction finishes then only we would commit
	// All or None concept
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
}
