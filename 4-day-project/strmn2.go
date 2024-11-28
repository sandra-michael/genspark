package main

import (
	"log"
	"strings"
)

// Function Custom type
type operation func(string) string

func main() {

	log.Println("function name :  trimSpace, return value: ", stringManipulation(trimSpace(), "          gfdngbk           "))
	log.Println("function name :  toUpper,  return value: ", stringManipulation(toUpper(), "sdfs"))
	log.Println("function name :  greetMe,  return value: ", stringManipulation(greetMe(), "Sandra"))
}

func stringManipulation(fn operation, input string) string {
	log.Println("----------------Execution Step 2 - Calling function which was returned----------------")
	return fn(input)
}

func trimSpace() operation {
	log.Println("---------------- Execution Step 1 - Within TrimSPace----------------")
	return func(input string) string {
		log.Println("----------------Execution Step 3 - within returned trim space----------------")
		return strings.TrimSpace(input)
	}
}

func toUpper() operation {
	log.Println("----------------Execution Step 1 - Within toUpper----------------")
	return func(input string) string {
		log.Println("----------------Execution Step 3 - within returned toUpper----------------")
		return strings.ToUpper(input)
	}

}

func greetMe() operation {
	log.Println("----------------Execution Step 1 - Within GreetMe----------------")
	return func(input string) string {
		log.Println("----------------Execution Step 3 - within returned greet----------------")
		return "Hello, " + input
	}
}
