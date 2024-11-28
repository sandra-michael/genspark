package mysql

import (
	"app/stores/models"
	"fmt"
)

type Conn struct {
	db string
}

func NewConn(con string) Conn {
	NewCon := Conn{con}
	return NewCon
}
func (c Conn) Create(u models.User) error {

	fmt.Println("Creating a user in ", c.db, " u : ", u)
	return nil

}

func (c Conn) CreateSimple(u string) error {
	if u != "" {
		fmt.Println("Creating a user in ", c.db, " u : ", u)
		return nil
	}
	return nil
}

func (c Conn) Update(name string) error {
	if name != "" {
		fmt.Println("UPdating a user in ", c.db, " name: ", name)
		return nil
	}
	return nil
}

func (c Conn) Delete(id int) error {
	if id != 0 {
		fmt.Println("deleting a user in ", c.db, " with id ", id)
		return nil
	}
	return nil
}
