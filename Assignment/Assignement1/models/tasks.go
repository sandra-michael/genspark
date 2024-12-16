package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conn struct {
	db *pgxpool.Pool
}

func NewConn() (*Conn, error) {
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
		return &Conn{}, err
	}

	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	config.MaxConns = 30

	config.HealthCheckPeriod = 5 * time.Minute

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return &Conn{}, err
	}
	return &Conn{db: db}, nil
}

func (c *Conn) Ping() {
	// pinging the connection if it is alive or not
	err := c.db.Ping(context.Background())
	if err != nil {
		panic(err)
	}
}

func (c *Conn) CreateTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        name VARCHAR(100),
		description VARCHAR(100),
    );`
	res, err := c.db.Exec(context.Background(), query)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("unable to create table: %w", err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())
	return nil

}

func (c *Conn) CreateTask(ctx context.Context, newTask NewTask) (int, error) {

	query := `INSERT INTO tasks (name,description,status) VALUES ($1,$2,$3) RETURNING id`
	var id int

	err := c.db.QueryRow(ctx, query, newTask.Name, newTask.Description, "NEW").Scan(&id)
	if err != nil {
		log.Fatalf("Unable to insert task: %v\n", err) // Log and terminate if user insertion fails
		return 0, fmt.Errorf("unable to insert task: %w", err)
	}
	fmt.Println("task inserted with id", id)
	return id, nil // Return the new user's ID
}

func (c *Conn) FetchTasks(ctx context.Context) ([]Task, error) {

	query := `SELECT * FROM tasks`

	// Execute the query to retrieve all tasks
	rows, err := c.db.Query(ctx, query)
	if err != nil {
		log.Fatalf("Unable to query users: %v\n", err) // Log and terminate if query fails
		return nil, fmt.Errorf("unable to fetch tasks: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed when done

	fmt.Println("tasks:")

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[Task])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return nil, fmt.Errorf("unable to fetch tasks: %w", err)
	}
	return tasks, nil

}

func (c *Conn) FetchTask(ctx context.Context, id int) (Task, error) {

	query := `SELECT * FROM tasks where id = $1`

	// Execute the query to retrieve all users
	rows, err := c.db.Query(ctx, query, id)
	if err != nil {
		fmt.Printf("Unable to query task: %v\n", err)
		return Task{}, fmt.Errorf("unable to fetch task: %w", err)
	}
	defer rows.Close() // Ensure the rows are closed when done

	fmt.Println("tasks:")

	task, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Task])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return Task{}, fmt.Errorf("unable to fetch task: %w", err)
	}
	return task, nil

}

func (c *Conn) UpdateTaskStatus(ctx context.Context, id int) error {
	selectQuery := `
		SELECT
			id, status
		FROM
			tasks
		WHERE
			id = $1
	`

	var task Task

	// Execute the query and scan the result into the book struct
	err := c.db.QueryRow(ctx, selectQuery, id).Scan(
		&task.ID,
		&task.Status,
	)
	if err != nil {
		return fmt.Errorf("unable to fetch task: %w", err)
	}

	updateStatus := "NEW"
	switch task.Status {
	case "NEW":
		updateStatus = "IN PROGRESS"
	case "IN PROGRESS":
		updateStatus = "DONE"
	case "DONE":
		return nil

	}
	// SQL query to update a user's email
	query := `UPDATE tasks SET status = $1 WHERE id = $2`

	_, err = c.db.Exec(context.Background(), query, updateStatus, id)
	if err != nil {
		log.Fatal("Unable to update task: %v\n", err) // Log and terminate if update fails
		return fmt.Errorf("unable to update task: %w", err)
	}
	fmt.Println("User status updated")
	return nil
}

func (c *Conn) UpdateTask(ctx context.Context, id int, updateTask UpdateTask) error {

	// Steps
	/*
		1. Add transactions
		2. Add validation to models.Book
		3. If validation fails then rollback the update and report some error to the user
	*/
	selectQuery := `
		SELECT
			id, name, description, status
		FROM
			tasks
		WHERE
			id = $1
	`

	var task Task

	// Execute the query and scan the result into the book struct
	err := c.db.QueryRow(ctx, selectQuery, id).Scan(
		&task.ID,
		&task.Name,
		&task.Description,
		&task.Status,
	)
	if err != nil {
		return fmt.Errorf("unable to fetch task: %w", err)
	}
	data, err := json.Marshal(updateTask)
	if err != nil {
		return fmt.Errorf("unable to marshal task: %w", err)
	}
	err = json.Unmarshal(data, &task)
	if err != nil {
		return fmt.Errorf("unable to unmarshal task: %w", err)
	}

	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()
	query := `
		UPDATE tasks
		SET name = $1, description = $2, status = $3
		WHERE id = $4
	`

	// Update the book based on its ID
	_, err = tx.Exec(ctx, query,
		task.Name, task.Description, task.Status, id,
	)

	if err != nil {
		return fmt.Errorf("unable to update task: %w", err)
	}

	fmt.Printf("Task with ID %d updated successfully\n", task.ID)
	return nil
}

func (c *Conn) DeleteTask(ctx context.Context, taskID int) error {
	// SQL query to delete a user by ID
	query := `DELETE FROM tasks WHERE id = $1`

	// Execute the query to delete the user
	_, err := c.db.Exec(ctx, query, taskID)
	if err != nil {
		fmt.Println("Unable to delete task: %v\n", err) // Log and terminate if deletion fails
		return fmt.Errorf("unable to delete tasks: %w", err)
	}
	fmt.Println("task deleted")
	return nil
}
