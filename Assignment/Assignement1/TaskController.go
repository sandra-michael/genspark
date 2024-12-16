package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var Db *pgxpool.Pool

func main() {
	fmt.Println("in main")
	fmt.Println("Creating Db")
	Db, err := CreateConnection()
	if err != nil {
		panic(err)
	}

	defer Db.Close()
	Ping(Db)
	//needs to return a response
	CreateTable(Db)

	mux := chi.NewRouter()

	mux.Route("/v1/tasks", func(r chi.Router) {
		//health check
		r.GET("/health", healthCheck)
		// POST    /api/v1/tasks                 # Create a new task
		r.POST("/", createTask)
		// GET     /api/v1/tasks/:id             # Get a specific task by ID
		r.GET("/{id}", func(w http.ResponseWriter, r *http.Request) {})
		// GET     /api/v1/tasks                 # Get all tasks (filter/sort optional)
		r.GET("/", func(w http.ResponseWriter, r *http.Request) {})
		// PUT     /api/v1/tasks/:id             # Update task details
		r.PUT("/{id}", func(w http.ResponseWriter, r *http.Request) {})
		// DELETE  /api/v1/tasks/:id             # Delete a task
		r.DELETE("/{id}", func(w http.ResponseWriter, r *http.Request) {})
		// PATCH   /api/v1/tasks/:id/status      # Update task status (e.g., In Progress, Done)
		r.PATCH("/{id}/status", func(w http.ResponseWriter, r *http.Request) {})
	})

	err := http.ListenAndServe(":8085", mux)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DOne")

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I am working and healthy"))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	//id := insertUser(context.Background(), db, "Jane", "<EMAIL>", 24)
	Name := r.URL.Query().Get("name")
	id, err := CreateTasks(r.Context(), Db, Name)
	if err != nil {

	}

	w.Write([]byte("Created Task with id : " + id))
}
