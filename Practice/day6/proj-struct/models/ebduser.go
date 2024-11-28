package models

import (
	"fmt"
	"proj-struct/database"
)

// create a struct for models, and make createUser as a method, and access the conf value using that
type EbdUser struct {
	name string
	database.Conf
}

func (user EbdUser) CreateUser(name string) {
	fmt.Println("creating the user", name)

	user.name = name
	//user.Conf = con

	//var con database.Conf
	user.Ping()

	fmt.Println("adding to db", user.Conf)
}

// func CreateUser(name string, c database.Conf) {
// 	fmt.Println("adding to db", c)
// 	fmt.Println("creating the user", name)
// }
