package main

import (
	"fmt"
	"proj-struct/database"
	"proj-struct/models"
)

// type user struct {
// 	name string
// }

func main() {

	// c.db = "mysql" // not allowed, db is not exported
	c := database.NewConf("mysql")
	fmt.Println(c)
	c.Ping()

	var s models.Service
	// s.Conf = c
	s = models.NewService(c)
	s.CreateUser("ajay")

	// var ebduser models.EbdUser
	// ebduser.CreateUser("someone")

}
