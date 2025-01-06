package products

import (
	"time"
)

// User struct represents the users table in the stores
type Product struct {
	ID          string    `json:"id"` // UUID
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	Category    string    `json:"category"`
	Stock       string    `json:"stock"`
	CreatedAt   time.Time `json:"created_at"` // Timestamp of creation
	UpdatedAt   time.Time `json:"updated_at"` // Timestamp of last update
}

// NewUser struct represents the data required when creating a new user
type NewProduct struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`        // Valid email required
	Description string `json:"description" validate:"required,min=2,max=100"` // Password must be at least 5 characters long
	Price       string `json:"price" validate:"required,min=2,max=100"`
	Category    string `json:"category" validate:"required,min=2,max=100"`
	Stock       string `json:"stock" validate:"required,min=2,max=100"`
}

// keeping it simple this is json which will be returned
type ProductOrder struct {
	ProductId string `json:"product_id"`
	PriceId   string `json:"price_id"`
	Stock     int    `json:"stock"`
}

//required: The roles field is mandatory.
//unique: Ensures there are no duplicate roles in the array (requires the validator's unique tag to be supported).
//dive: Applies validation rules to each individual element of the slice.
//oneof=user admin: Restricts each role value to either user or admin.

type ProductDetail struct {
	ID          string `json:"id"` // UUID
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
}

// Product Order request
type ProductOrdersRequest struct {
	ProductIDs []string `json:"productIds" binding:"required"`
}

// ProductUpdateRequest represents the fields that can be updated
type ProductUpdateRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`        // Valid email required
	Description string `json:"description,omitempty" validate:"omitempty,min=2,max=100"` // Password must be at least 5 characters long
	//TODO allow when stripe price update is sorted
	//Price       string `json:"price,omitempty" validate:"omitempty,min=2,max=100"`
	Category string `json:"category,omitempty" validate:"omitempty,min=2,max=100"`
	Stock    string `json:"stock,omitempty" validate:"omitempty,min=2,max=100"`
}

/*
	/*
		//------------------------------------------------------//
		//   Adding Cart Service Structs
		//------------------------------------------------------//
*/

type Cart struct {
	ID        string    `json:"id"`         // Maps to UUID PRIMARY KEY
	ProductID string    `json:"product_id"` // Maps to UUID, foreign key to products table
	UserID    string    `json:"user_id"`    // Maps to UUID, identifies the user
	OrderID   string    `json:"order_id"`   // Maps to UUID, identifies the order
	Quantity  int       `json:"quantity"`   // Maps to INTEGER, must be >= 1
	Status    string    `json:"status"`     // Maps to ENUM 'status_enum', defaults to 'inprogress'
	CreatedAt time.Time `json:"created_at"` // Maps to TIMESTAMP
	UpdatedAt time.Time `json:"updated_at"` // Maps to TIMESTAMP
}

type StatusEnum string

const (
	StatusInProgress StatusEnum = "inprogress"
	StatusPending    StatusEnum = "pending"
	StatusCompleted  StatusEnum = "completed"
)

// NewUser struct represents the data required when creating a new line item inside the cart
type NewCartLine struct {
	ProductID string `json:"product_id" validate:"required,min=2,max=100"`
	Quantity  int    `json:"quantity" validate:"required,min=1,max=100"`
}

type LineItem struct {
	ProductID string `json:"productId" `
	Quantity  int    `json:"quantity"`
}

type CartReturn struct {
	OrderId   string
	LineItems []LineItem
}

type OrderRequest struct {
	LineItems []LineItem `json:"lineItems" `
}

type FetchCartResponse struct {
	LineItems []LineItem `json:"lineItems" binding:"required"`
}

type CartDetails struct {
	ID        string
	OrderId   string
	ProductID string
	Quantity  int
}
