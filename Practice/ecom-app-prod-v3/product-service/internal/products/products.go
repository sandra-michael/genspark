package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Conf struct {
	db *sql.DB
}

func NewConf(db *sql.DB) (*Conf, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	return &Conf{db: db}, nil
}

func (c *Conf) InsertProduct(ctx context.Context, newProduct NewProduct) (Product, error) {

	id := uuid.NewString()

	// Get the current UTC time for `createdAt` and `updatedAt` timestamps for the new user.
	createdAt := time.Now().UTC()
	updatedAt := time.Now().UTC()

	var prod Product

	// Use a transaction to ensure atomicity of the database operation.

	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// SQL query to insert a new user into the "users" table.
		// The `RETURNING` clause retrieves the inserted user's data after the operation.
		query := `
      INSERT INTO products
      (id, name, description, price, category, stock, created_at,updated_at)
      VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
      RETURNING id, name, created_at, updated_at
      `
		// Execute the `INSERT` query within the transaction to add the new user.
		// `QueryRowContext` executes the query and scans the resulting row into the `user` struct.
		err := tx.QueryRowContext(ctx, query, id, newProduct.Name, newProduct.Description, newProduct.Price, newProduct.Category, newProduct.Stock, createdAt, updatedAt).
			Scan(&prod.ID, &prod.Name, &prod.CreatedAt, &prod.UpdatedAt)
		if err != nil {
			// Return an error if the query execution or scan fails.
			return fmt.Errorf("failed to insert user: %w", err)
		}

		// If the query is successful, return nil to indicate no errors.
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return Product{}, fmt.Errorf("failed to insert user: %w", err)
	}

	return prod, nil

}

// withTx is a helper function that simplifies the usage of SQL transactions.
// It begins a transaction using the provided context (`ctx`), executes the given function (`fn`),
// and handles commit or rollback based on the success or failure of the function.
func (c *Conf) withTx(ctx context.Context, fn func(*sql.Tx) error) error {
	// Start a new transaction using the context.
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err) // Return an error if the transaction cannot be started.
	}

	// Execute the provided function (`fn`) within the transaction.
	if err := fn(tx); err != nil {
		// If the function returns an error, attempt to roll back the transaction.
		er := tx.Rollback()
		if er != nil && !errors.Is(err, sql.ErrTxDone) {
			// If rollback also fails (and it's not because the transaction is already done),
			// return an error indicating the failure to roll back.
			return fmt.Errorf("failed to rollback withTx: %w", err)
		}
		// Return the original error from the function execution.
		return fmt.Errorf("failed to execute withTx: %w", err)
	}

	// If no errors occur, commit the transaction to apply the changes.
	err = tx.Commit()
	if err != nil {
		// Return an error if the transaction commit fails.
		return fmt.Errorf("failed to commit withTx: %w", err)
	}

	// Return nil if the function executes successfully and the transaction is committed.
	return nil
}
