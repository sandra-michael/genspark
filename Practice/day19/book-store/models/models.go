package models

type Book struct {
	ID          int     `json:"id"`    // Unique identifier (primary key in SQL)
	Title       string  `json:"title"` // Title of the book
	AuthorName  string  `json:"author_Name"`
	AuthorEmail string  `json:"author_email"`
	Description string  `json:"description"` // Description of the book
	Category    string  `json:"category"`    // Book category (e.g., Fiction, Biography)
	Price       float64 `json:"price"`       // Price of the book
	Stock       int     `json:"-"`           // Number of copies in stock
}

type NewBook struct {
	Title       string  `json:"title"` // Title of the book
	AuthorName  string  `json:"author_Name"`
	AuthorEmail string  `json:"author_email"`
	Description string  `json:"description"` // Description of the book
	Category    string  `json:"category"`    // Book category (e.g., Fiction, Biography)
	Price       float64 `json:"price"`       // Price of the book
	Stock       int     `json:"stock"`       // Number of copies in stock
}
