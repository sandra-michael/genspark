package products

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

func (c *Conf) DecrementStock(ctx context.Context, productId string, stock int) error {
	updatedAt := time.Now().UTC() // Current timestamp

	// Use a transaction to ensure consistency
	err := c.withTx(ctx, func(tx *sql.Tx) error {

		// Step 3: Perform the update since the `updated_at` condition is met
		queryUpdate := `
		UPDATE products
		SET stock = stock - $2,  updated_at = $3
		WHERE id = $1 AND stock > 0;

		`

		res, err := tx.ExecContext(ctx, queryUpdate, productId, stock, updatedAt)
		if err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}

		num, err := res.RowsAffected()
		if num == 0 || err != nil {
			return fmt.Errorf("failed to update order: %w", err)
		}

		// Successfully updated the order
		return nil
	})

	if err != nil {
		// Return the error, if any
		return err
	}

	// Return nil if the update is successful or skipped gracefully
	return nil
}

func (c *Conf) FetchAllProducts(ctx context.Context) ([]ProductDetail, error) {

	// An album slice to hold data from returned rows.
	var products []ProductDetail
	// Use a transaction to ensure consistency
	err := c.withTx(ctx, func(tx *sql.Tx) error {

		// Step 3: Perform the update since the `updated_at` condition is met
		queryFetch := `
		select id,name,description ,price,stock
		from products;
		`

		rows, err := tx.Query(queryFetch)
		if err != nil {
			return fmt.Errorf("failed to fetch product: %w", err)
		}

		defer rows.Close()

		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var prod ProductDetail
			if err := rows.Scan(&prod.ID, &prod.Name, &prod.Description,
				&prod.Price, &prod.Stock); err != nil {
				return err
			}
			products = append(products, prod)
		}
		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
	}

	return products, nil

}

func (c *Conf) UpdateProduct(ctx context.Context, productId string, req ProductUpdateRequest) error {

	// Build the update query dynamically
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		// SQL query to insert a new user into the "users" table.
		// The `RETURNING` clause retrieves the inserted user's data after the operation.
		query, args, err := buildUpdateQuery("products", "id", productId, req)
		if err != nil {
			return err
		}

		// Execute the update query
		queryWithReturning := query + " RETURNING id, name, description, price, stock"
		_, err = tx.ExecContext(ctx, queryWithReturning, args...)
		if err != nil {
			return err
		}
		// If the query is successful, return nil to indicate no errors.
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	return nil

}

// TODO ADD TRANSACTIONS
// buildUpdateQuery dynamically constructs the SQL update query
// buildUpdateQuery dynamically constructs the SQL update query
func buildUpdateQuery(tableName, idColumn, idValue string, req ProductUpdateRequest) (string, []interface{}, error) {
	setClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	updatedAt := time.Now().UTC() // Current timestamp
	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argIndex))
	args = append(args, updatedAt)
	argIndex++

	// Dynamically add columns to update based on non-nil fields
	if req.Name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", argIndex))
		args = append(args, req.Name)
		argIndex++
	}
	if req.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, req.Description)
		argIndex++
	}
	if req.Category != "" {
		setClauses = append(setClauses, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, req.Category)
		argIndex++
	}
	if req.Stock != "" {
		setClauses = append(setClauses, fmt.Sprintf("stock = $%d", argIndex))
		args = append(args, req.Stock)
		argIndex++
	}

	// If no fields to update, return an error
	if len(setClauses) == 0 {
		return "", nil, errors.New("no fields to update")
	}

	// Construct the final query
	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s = $%d", tableName,
		joinStrings(setClauses, ", "), idColumn, argIndex)
	args = append(args, idValue)

	return query, args, nil
}

// joinStrings is a utility to join strings with a delimiter
func joinStrings(elements []string, delimiter string) string {
	return fmt.Sprintf(strings.Join(elements, delimiter))
}

//https://docs.stripe.com/products-prices/manage-prices?dashboard-or-api=api#archive-price
//To delete a product
//need to archive price
//then delete the product
//TODO add delete api for procuxt
