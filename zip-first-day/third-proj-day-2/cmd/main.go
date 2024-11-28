package main

import (
	"fmt"
	"third-proj-day-2/auth"
)

func main() {
	setup()
	auth.Authenticate()
	fmt.Println("end of the main")
}
