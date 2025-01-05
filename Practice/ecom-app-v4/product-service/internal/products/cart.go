package products

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (c *Conf) InsertOrUpdateCart(ctx context.Context, userId string, lineItem NewCartLine) error {

	id := uuid.NewString()

	// Use a transaction to ensure consistency
	err := c.withTx(ctx, func(tx *sql.Tx) error {
		var existingOrderID string
		// Check if user exists with a pending status
		err := tx.QueryRow(`
		SELECT order_id 
		FROM cart 
		WHERE user_id = $1 AND status = 'inprogress' 
		LIMIT 1
	`, userId).Scan(&existingOrderID)

		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("failed to query existing order ID: %w", err)
		}

		if existingOrderID != "" {
			// Check if product exists for the user with pending status
			var currentQuantity int
			err = tx.QueryRow(`
			SELECT quantity 
			FROM cart 
			WHERE user_id = $1 AND product_id = $2 AND status = 'inprogress'
		`, userId, lineItem.ProductID).Scan(&currentQuantity)

			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("failed to query existing product: %w", err)
			}

			if currentQuantity > 0 {
				// Update the quantity for the existing product
				_, err = tx.Exec(`
				UPDATE cart 
				SET quantity = quantity + $1, updated_at = $2 
				WHERE user_id = $3 AND product_id = $4 AND status = 'inprogress'
			`, lineItem.Quantity, time.Now().UTC(), userId, lineItem.ProductID)

				if err != nil {
					return fmt.Errorf("failed to update product quantity: %w", err)
				}
			} else {
				// Insert a new row with the existing order ID
				_, err = tx.Exec(`
				INSERT INTO cart (id, product_id, user_id, order_id, quantity, status, created_at, updated_at) 
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			`, id, lineItem.ProductID, userId, existingOrderID, lineItem.Quantity, StatusInProgress, time.Now().UTC(), time.Now().UTC())

				if err != nil {
					return fmt.Errorf("failed to insert new product: %w", err)
				}
			}
		} else {
			// Insert a new row with a new order ID
			newOrderId := uuid.NewString()

			_, err = tx.Exec(`
			INSERT INTO cart (id, product_id, user_id, order_id, quantity, status, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`, id, lineItem.ProductID, userId, newOrderId, lineItem.Quantity, StatusInProgress, time.Now().UTC(), time.Now().UTC())

			if err != nil {
				return fmt.Errorf("failed to insert new row with new order ID: %w", err)
			}
		}
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}

	return nil

}

func (c *Conf) FetchCartItems(ctx context.Context, userId string, status StatusEnum) (CartReturn, error) {

	// An album slice to hold data from returned rows.
	var cartLines []LineItem

	var ret CartReturn
	var order_id string
	// Use a transaction to ensure consistency
	err := c.withTx(ctx, func(tx *sql.Tx) error {

		// Step 3: Perform the update since the `updated_at` condition is met
		queryFetch := `
		select order_id,product_id, quantity
		from cart
		WHERE user_id = $1 AND status = $2;
		`

		rows, err := tx.Query(queryFetch, userId, status)
		if err != nil {
			return fmt.Errorf("failed to fetch cart items: %w", err)
		}

		defer rows.Close()

		// Loop through rows, using Scan to assign column data to struct fields.
		for rows.Next() {
			var line LineItem

			if err := rows.Scan(&order_id, &line.ProductID, &line.Quantity); err != nil {
				return err
			}
			cartLines = append(cartLines, line)
		}
		if err = rows.Err(); err != nil {
			return err
		}
		return nil
	})

	// If the transaction or insertion fails, return an error.
	if err != nil {
		return CartReturn{}, fmt.Errorf("failed to fetch products: %w", err)
	}
	ret.OrderId = order_id
	ret.LineItems = cartLines

	return ret, nil

}

func (c *Conf) UpdateCartStatusFromInProgressToPending(ctx context.Context, userId string) error {

	updatedAt := time.Now().UTC() // Current timestamp

	// Use a transaction to ensure consistency
	err := c.withTx(ctx, func(tx *sql.Tx) error {

		// Step 3: Perform the update since the `updated_at` condition is met
		queryUpdate := `
				UPDATE cart 
				SET status = $1, updated_at = $2 
				WHERE user_id = $3 AND status = 'inprogress'
			`

		res, err := tx.ExecContext(ctx, queryUpdate, StatusPending, updatedAt, userId)
		if err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
		}

		num, err := res.RowsAffected()
		if num == 0 || err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
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

func (c *Conf) UpdateCartStatusForOrderId(ctx context.Context, orderId string) error {

	updatedAt := time.Now().UTC() // Current timestamp

	// Use a transaction to ensure consistency
	err := c.withTx(ctx, func(tx *sql.Tx) error {

		// Step 3: Perform the update since the `updated_at` condition is met
		queryUpdate := `
				UPDATE cart 
				SET status = $1, updated_at = $2 
				WHERE order_id = $3
			`

		res, err := tx.ExecContext(ctx, queryUpdate, StatusCompleted, updatedAt, orderId)
		if err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
		}

		num, err := res.RowsAffected()
		if num == 0 || err != nil {
			return fmt.Errorf("failed to update cart status: %w", err)
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
