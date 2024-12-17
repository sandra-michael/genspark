package handlers

import (
	"Assignment1/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	//here we are adding service layer and are decoupling the use of *models.COnn
	//this helps with testing
	c        models.Service
	validate *validator.Validate
}

func NewConf(c models.Service, validate *validator.Validate) *Handler {
	return &Handler{c: c, validate: validate}
}

func SetUpTaskRoute(conn models.Service) *chi.Mux {
	h := NewConf(conn, validator.New())
	mux := chi.NewRouter()

	mux.Route("/v1/tasks", func(r chi.Router) {
		//health check
		r.Get("/health", h.healthCheck)
		// // POST    /api/v1/tasks                 # Create a new task
		r.Post("/", h.createTask)
		// // GET     /api/v1/tasks/:id             # Get a specific task by ID
		r.Get("/{id}", h.fetchTask)
		// // GET     /api/v1/tasks                 # Get all tasks (filter/sort optional)
		r.Get("/", h.fetchTasks)
		// // PUT     /api/v1/tasks/:id             # Update task details
		r.Put("/{id}", h.updateTask)
		// // DELETE  /api/v1/tasks/:id             # Delete a task
		r.Delete("/{id}", h.deleteTask)
		// // PATCH   /api/v1/tasks/:id/status      # Update task status (e.g., In Progress, Done)
		r.Patch("/{id}/status", h.UpdateTaskStatus)
	})
	return mux
}
