package main

import (
	"log"
	"strings"
)

// q3. Create a function named as StringManipulation.
//     StringManipulation accepts a function and string type as an argument, and it returns string value
//     Possible Functions that it can accept:- trimSpace, toUpper, greet

//     Create 3 functions trimSpace, toUpper, greet
//     TrimSpace:- TrimSpace returns a string, with all leading and trailing white space removed, as defined by Unicode.
//     ToUpper:- ToUpper returns string with all Unicode letters mapped to their upper case.
//     Greet: - It takes a name as input, add hello as greeting and return the greeting
//     Hint: use strings package for TrimSpace and ToUpper functionalities

func main() {
	// stringManipulation(trimSpace, "          gfdngbk           ")
	// stringManipulation(toUpper, "sdfs")
	// stringManipulation(greetMe, "Sandra")
	log.Println("function name :  trimSpace, return value: ", stringManipulation(trimSpace, "          gfdngbk           "))
	log.Println("function name :  toUpper,  return value: ", stringManipulation(toUpper, "sdfs"))
	log.Println("function name :  greetMe,  return value: ", stringManipulation(greetMe, "Sandra"))
}

func stringManipulation(fn func(string) string, input string) string {
	//log.Println("function name : ", fn, " return value: ", fn(input))
	return fn(input)
}

func trimSpace(input string) string {
	return strings.TrimSpace(input)
}

func toUpper(input string) string {
	return strings.ToUpper(input)
}

func greetMe(input string) string {
	return "Hello, " + input
}
