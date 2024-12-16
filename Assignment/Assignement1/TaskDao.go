package main

import (
	"Assignment1/models"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"log"
	"time"
)

// func main() {
// 	fmt.Println("in main")
// 	fmt.Println("Creating Db")
// 	Db, err := CreateConnection()
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer Db.Close()
// 	Ping(Db)
// 	//needs to return a response
// 	CreateTable(Db)
// }

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

	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

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

func CreateTable(db *pgxpool.Pool) {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100)
    );`
	res, err := db.Exec(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())

}

func CreateTasks(ctx context.Context, db *pgxpool.Pool, name string) (int, error) {

	query := `INSERT INTO tasks (name) VALUES ($1) RETURNING id`
	var id int

	err := db.QueryRow(ctx, query, name).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to insert task: %v\n", err) // Log and terminate if user insertion fails
		return 0, err
	}
	fmt.Println("task inserted with id", id)
	return id, nil // Return the new user's ID
}

func FetchTasks(ctx context.Context, db *pgxpool.Pool) ([]models.Task, error) {

	query := `SELECT id, name FROM tasks`

	// Execute the query to retrieve all users
	rows, err := db.Query(ctx, query)
	if err != nil {
		log.Fatalf("Unable to query users: %v\n", err) // Log and terminate if query fails
	}
	defer rows.Close() // Ensure the rows are closed when done

	fmt.Println("tasks:")

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Task])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return nil, err
	}
	return tasks, nil

}

func FetchTask(ctx context.Context, db *pgxpool.Pool, id int) (models.Task, error) {

	query := `SELECT id, name FROM tasks where id = $1`

	// Execute the query to retrieve all users
	rows, err := db.Query(ctx, query)
	if err != nil {
		fmt.Printf("Unable to query task: %v\n", err)
		return models.Task{}, err
	}
	defer rows.Close() // Ensure the rows are closed when done

	fmt.Println("tasks:")

	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[models.Task])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return models.Task{}, err
	}
	return task, nil

}

func UpdateTask(db *pgxpool.Pool, taskID int, name string) error {
	// SQL query to update a user's email
	query := `UPDATE tasks SET name = $1 WHERE id = $2`

	_, err := db.Exec(context.Background(), query, name, taskID)
	if err != nil {
		log.Fatal("Unable to update task: %v\n", err) // Log and terminate if update fails
		return err
	}
	fmt.Println("User email updated")
	return nil
}

func DeleteTask(ctx context.Context, taskID int, db *pgxpool.Pool) {
	// SQL query to delete a user by ID
	query := `DELETE FROM task WHERE id = $1`

	// Execute the query to delete the user
	_, err := db.Exec(ctx, query, taskID)
	if err != nil {
		log.Fatalf("Unable to delete task: %v\n", err) // Log and terminate if deletion fails
	}
	fmt.Println("task deleted")
}
