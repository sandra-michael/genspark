package models

import (
	"fmt"
	"proj-struct/database"
)

// create a struct for models, and make createUser as a method, and access the conf value using that
type User struct {
	name string
}

type Service struct {
	database.Conf
}

func (s *Service) CreateUser(name string) {
	fmt.Println("creating the user", name)

	var user User
	user.name = name

	//var con database.Conf
	//s.Ping()

	fmt.Println("adding to db", s.Conf)
}

func NewService(c database.Conf) Service {
	service := Service{Conf: c}
	return service
}

func CreateUser(name string, c database.Conf) {
	fmt.Println("adding to db", c)
	fmt.Println("creating the user", name)
}
