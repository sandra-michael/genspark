package main

import (
	"fmt"
	"log"
)

type user struct {
	//io.Writer // panic: runtime error: invalid memory address or nil pointer dereference // if directly embedded
	name  string
	email string
}

// hence we go with the concrete approach
func (u user) Write(p []byte) (n int, err error) {
	fmt.Printf("sending a notification to %s %s %s", u.name, u.email, string(p))
	return len(p), nil

}
func main() {
	u := user{name: "raj", email: "raj@email.com"}
	l := log.New(u, "log: ", log.LstdFlags)
	//name := "newname"
	//l := log.New(name, "log: ", log.LstdFlags)

	l.Println("Hello, log")
}
