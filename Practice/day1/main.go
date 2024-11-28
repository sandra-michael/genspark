package main

import (
	"first-proj-day-2/stringops"
	"fmt"
)

func main() {
	fmt.Println("your in main ")
	s1, s2 := "test1", "test2"
	fmt.Println(stringops.ReverseAndUppercase(s1, s2))
}
