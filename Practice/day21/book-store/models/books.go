package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.23.0"
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

	return &Conn{db: db}, nil
}

// func main() {
// 	conn, err := NewConn()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}

// 	err = conn.CreateBookTable(context.Background())
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// }

func (c *Conn) CreateBookTable(ctx context.Context) error {
	query := `
    CREATE TABLE IF NOT EXISTS books (
        id SERIAL PRIMARY KEY,
        title VARCHAR(100),
		author_name VARCHAR(100),	
		author_email VARCHAR(100),
		description VARCHAR(100),
		category VARCHAR(100),
		price INT,
		stock INT

    );`
	res, err := c.db.Exec(context.Background(), query)
	if err != nil {
		//log.Fatal(err)
		return fmt.Errorf("unable to insert book: %w", err)
	}
	fmt.Printf("rows affected: %d\n", res.RowsAffected())
	return nil
}

func (c *Conn) InsertBook(ctx context.Context, newBook NewBook) (Book, error) {
	tracer := otel.Tracer("book-micro")
	ctx, span := tracer.Start(ctx, "InsertBook")
	defer span.End()

	traceId := span.SpanContext().TraceID().String()

	query := `
		INSERT INTO books (
		                   title, author_name,author_email, 
		                   description, category, 
		                   price, stock
		                  
		                   )
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	var id int
	err := c.db.QueryRow(
		ctx, query, newBook.Title, newBook.AuthorName,
		newBook.AuthorEmail, newBook.Description, newBook.Category,
		newBook.Price, newBook.Stock,
	).Scan(&id)

	if err != nil {
		//log.Println(err)
		fmt.Println(err)
		// Handle and record any errors in the span
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // HTTP 400
		//span.SetAttributes(attribute.String("user_id", userId))                          // Attach user ID
		span.SetAttributes(attribute.String("traceId", traceId))
		span.AddEvent("UNABLE TO INSERT BOOK in db") // Record event in tracing span// Attach trace ID
		span.SetStatus(codes.Error, err.Error())
		return Book{}, fmt.Errorf("unable to insert book: %w", err)
	}

	b := Book{
		ID:          id,
		Title:       newBook.Title,
		AuthorName:  newBook.AuthorName,
		AuthorEmail: newBook.AuthorEmail,
		Description: newBook.Description,
		Category:    newBook.Category,
		Price:       newBook.Price,
		Stock:       newBook.Stock,
	}
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK)) // HTTP 200

	return b, nil

}
