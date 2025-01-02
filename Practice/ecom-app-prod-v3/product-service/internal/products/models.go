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
	PriceId string `json:"price_id"`
	Stock   int `json:"stock"`
}

//required: The roles field is mandatory.
//unique: Ensures there are no duplicate roles in the array (requires the validator's unique tag to be supported).
//dive: Applies validation rules to each individual element of the slice.
//oneof=user admin: Restricts each role value to either user or admin.
