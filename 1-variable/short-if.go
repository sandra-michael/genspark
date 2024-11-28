package main

import (
	"fmt"
)

func main() {

	if i := 10; i < 10 {
		fmt.Println("i is less than 10")
	}
	//a, _ := f() // assume f() return two values, the first value would be stored in a var
	// , and the second value would be ignored because of _(underscore)

	// use short if
	// call the println func and check if it has written no values over the stdout
	// if no values are written return otherwise print the values

	if a, _ := fmt.Println(); a == 1 { //1 because we have a new line character
		fmt.Println(a)
	}

	if a, error := fmt.Println(); a > 0 {
		fmt.Println(error)
	}
}
