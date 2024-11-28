package auth

import (
	"fmt"
	"second-proj-day-2/user"
)

func Authenticate() {
	fmt.Println(" authenticating user")
}

func Name() {
	fmt.Println(" authenticating user", user.GlobalName)
}
