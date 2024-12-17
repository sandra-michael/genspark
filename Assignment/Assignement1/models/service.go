package models

import "context"

//will use this service in handler to decouple the direct use of connection
//running go generate
//go generate

// We write go generate to run the below line, it would generate mock implementation of interface
// run go generate command from the current directory

// flags
// - source - fileName
// - destination - destination for generated mocks
// - package - package name for mock

//go:generate mockgen -source service.go -destination mockmodels/service_mock.go -package mockmodels
type Service interface {
	Ping()
	CreateTable() error
	CreateTask(ctx context.Context, newTask NewTask) (int, error)
	FetchTasks(ctx context.Context) ([]Task, error)
	FetchTask(ctx context.Context, id int) (Task, error)
	UpdateTaskStatus(ctx context.Context, id int) error
	UpdateTask(ctx context.Context, id int, updateTask UpdateTask) error
	DeleteTask(ctx context.Context, taskID int) error
}
