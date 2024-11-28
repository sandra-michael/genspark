package database

import (
	"fmt"
	"log"
)

type Conf struct {
	db string // not exporting the field, so no one can change the values outside this package
}

// if we are not exporting the field, then we can't set the values to it directly from outside,
// we need to create a function that initialize the struct
// New functions are used to initialize struct with some config values,

func NewConf(conn string) Conf {
	if conn == "" {
		// avoid in production until and unless you want your app to stop working
		// this will crash the program
		log.Fatal("empty connection string")
	}
	// try to open the connection, and if it is successful, return the connection
	return Conf{db: conn}
}

func (c Conf) Ping() {
	if c.db == "" {
		log.Fatal("can't ping the database, connection string is empty")
	}
	fmt.Println("pinging the database", c.db)

}
