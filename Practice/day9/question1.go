package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// q1. Create a function that converts string to an integer
//     if any alphabets are passed, wrap strconv error and ErrStringValue error (create ErrStringValue error)

//     ErrStringValue contains a message that 'value is of string type' and return the wrapped errors
//     otherwise return the original error

//     use the regex to check if value is of string type or not
//     Hint: regexp.MatchString(`^[a-zA-Z]`, s)
//     fmt.Errorf("%w %w") // to wrap error

//     In main function check if ErrStringValue error was wrapped in the chain or not
//     If yes, log a message 'value must be of int type not string' and log original error message alongside as well

var ErrStringValue error = errors.New("value is of string type")

func main() {
	ipVal := "test" // throws wrapped error along with invalid syntax
	//ipVal := "1234456"
	//ipVal := "1234456e2" // throws invalid syntax - ErrSyntax

	//throws an error but is not of string value
	//ipVal := "12344569897890090890097878767e2" // throws value out of range -  ErrRange
	opVal, errs := strToInt(ipVal)

	if errors.Is(errs, ErrStringValue) {
		fmt.Printf("value must be of int type not string, StackTrace : %v", errs)
		return
	}

	if errs != nil {
		fmt.Println(errs)
		return

	}
	fmt.Println("Converted integers val ", opVal, "of string ", ipVal)
}

func strToInt(s string) (int64, error) {

	i, err := strconv.ParseInt(s, 10, 32)

	if err != nil {
		ok, interr := regexp.MatchString(`^[a-zA-Z]`, s)
		if interr != nil {
			return 0, fmt.Errorf("%w %w", interr, err)
		}
		if ok {
			return 0, fmt.Errorf("%w %w", ErrStringValue, err)
		}
		return 0, err

	}
	return i, nil

}
