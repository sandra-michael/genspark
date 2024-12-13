package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

//https://hub.docker.com/_/postgres

//orm
//https://gorm.io/docs/query.html

// driver
//https://github.com/jackc/pgx

// whatever module we download using go get, is stored inside our gopath
// go env GOPATH

// go get moduleName (to get an external lib)
// github.com/jackc/pgx/v5 (don't forget to include the version number if there is a major version in the module name)

// automatically resolves all the dependencies required for the project
// go mod tidy // first command to run when importing any project

func CreateConnection() (*pgxpool.Pool, error) {

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

	// ParseConfig takes the connection string to generate a config
	config, err := pgxpool.ParseConfig(psqlInfo)
	if err != nil {
		return nil, err
	}

	// MinConns is the minimum number of connections kept open by the pool.
	// The pool will not proactively create this many connections, but once this many have been established,
	// it will not close idle connections unless the total number exceeds MaxConns.
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	// MaxConns is the maximum number of connections that can be opened to PostgreSQL.
	// This limit can be used to prevent overwhelming the PostgreSQL server with too many concurrent connections.
	config.MaxConns = 30

	config.HealthCheckPeriod = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Ping(db *pgxpool.Pool) {
	// pinging the connection if it is alive or not
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
}

// Three methods to execute queries on the database
// Exec -> when query does not return anything
// QueryRow -> returns exactly one row
// Query -> returns many rows

func main() {
	db, err := CreateConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	Ping(db)

	CreateTable(db)
	//id := insertUser(context.Background(), db, "John", "<EMAIL>", 25)
	//updateUserEmail(db, id, "john@example.com")
	id := insertUser(context.Background(), db, "Jane", "<EMAIL>", 24)
	updateUserEmail(db, id, "jane@example.com")
	id = insertUser(context.Background(), db, "Jill", "<EMAIL>", 23)
	updateUserEmail(db, id, "jill@example.com")

}

func CreateTable(db *pgxpool.Pool) {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
        email VARCHAR(100) UNIQUE NOT NULL,
        age INT
    );`
	res, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())

}

// Create two function one to insert one record and another to update the record

// insertUser inserts a new user into the users table and returns the new user's ID.
func insertUser(ctx context.Context, db *pgxpool.Pool, name, email string, age int) int {
	// SQL query to insert a user and return the new user's ID
	// don't hardcode the values, or use the string in construction, sql injection can happen
	query := `INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id`
	var id int

	// Execute the query to insert the user and get the new user's ID
	//QueryRow returns one row as output
	err := db.QueryRow(ctx, query, name, email, age).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to insert user: %v\n", err) // Log and terminate if user insertion fails
	}
	fmt.Println("User inserted with id", id)
	return id // Return the new user's ID
}

// updateUserEmail updates a user's email based on their ID.
func updateUserEmail(db *pgxpool.Pool, userID int, newEmail string) {
	// SQL query to update a user's email
	query := `UPDATE users SET email = $1 WHERE id = $2`

	// prone to SQL injection, we should not construct queries using strings
	//query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", "userInput")

	// Execute the query to update the user's email
	_, err := db.Exec(context.Background(), query, newEmail, userID)
	if err != nil {
		log.Fatalf("Unable to update user: %v\n", err) // Log and terminate if update fails
	}
	fmt.Println("User email updated")
}

func getAllUsers(ctx context.Context, db *pgxpool.Pool) {
	// user query method
	// run a for loop on rows.Next()
	// inside the loop scan values using rows.Scan
	// print things inside the loop

	// SQL query to retrieve all users
	query := `SELECT id, name, email, age FROM users`

	// Execute the query to retrieve all users
	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Fatalf("Unable to query users: %v\n", err) // Log and terminate if query fails
	}
	defer rows.Close() // Ensure the rows are closed when done

	fmt.Println("Users:")

	// this loop would run until there are rows to scan
	for rows.Next() {
		var id, age int
		var name, email string

		// Scan each row into variables
		err := rows.Scan(&id, &name, &email, &age)
		if err != nil {
			log.Printf("Unable to scan row: %v\n", err) // Log
			continue
		}
		fmt.Printf("ID: %d, Name: %s, Email: %s, Age: %d\n", id, name, email, age) // Print user details
	}

}

// deleteUser deletes a user from the users table based on their ID.
func deleteUser(ctx context.Context, userID int, db *pgxpool.Pool) {
	// SQL query to delete a user by ID
	query := `DELETE FROM users WHERE id = $1`

	// Execute the query to delete the user
	_, err := db.Exec(ctx, query, userID)
	if err != nil {
		log.Fatalf("Unable to delete user: %v\n", err) // Log and terminate if deletion fails
	}
	fmt.Println("User deleted")
}
