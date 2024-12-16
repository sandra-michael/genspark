package models

type Task struct {
	ID          int    `json:"id"` // Unique identifier (primary key in SQL)
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	// createdOn   time.Time `json:"created_on"`
	// modifiedOn  time.Time    `json:"modified_on"`
}

type NewTask struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description" validate:"required,min=3,max=500"`
}

type UpdateTask struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
