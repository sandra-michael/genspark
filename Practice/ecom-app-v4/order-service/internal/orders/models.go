package orders

import "time"

// Order represents an order entity in the database
type Order struct {
	ID                  int64     `json:"id"`                    // Auto-incrementing ID
	UserID              string    `json:"user_id"`               // UUID of the user placing the order
	ProductID           string    `json:"product_id"`            // UUID of the product
	Status              string    `json:"status"`                // Order status: pending, paid, or canceled
	StripeTransactionID string    `json:"stripe_transaction_id"` // Stripe transaction ID
	TotalPrice          int64     `json:"total_price"`           // Total price of the order in cents
	CreatedAt           time.Time `json:"created_at"`            // When the order was created
	UpdatedAt           time.Time `json:"updated_at"`            // When the order was last updated
}
