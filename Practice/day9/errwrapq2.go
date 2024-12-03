package main

import (
	"errors"
	"fmt"
	"runtime/debug"
)

// q2. Create 3 functions f1, f2, f3
//
//	f1() call f2(), f2() call f3()
//	each layer would return the error, wrap the error from each layer
//	print stack trace using debug.Stack to get a complete stack trace
var ErrFunc1 = errors.New("Error from function 1")
var ErrFunc2 = errors.New("Error from function 2")
var ErrFunc3 = errors.New("Error from function 3")

func main() {
	_, err := f1()

	if err != nil {
		fmt.Println("There was an error during f1 function call ")
		fmt.Println("error ", err)
		//debug.Stack() see's who called it
		//fmt.Println("stacktrace : ", string(debug.Stack()))
		if errors.Is(err, ErrFunc3) {
			fmt.Printf("value must be of int type not string, StackTrace : %v", err)
			return
		}
		return
	}
}

func f1() (int, error) {
	_, err := f2()

	if err != nil {
		return 0, fmt.Errorf("%w %w", ErrFunc1, err)
	}
	return 0, nil
}

func f2() (int, error) {
	_, err := f3()

	if err != nil {
		//if we use %v instead of %w this will not work in errors.Is(err, ErrFunc3)
		return 0, fmt.Errorf("%w %w", ErrFunc2, err)
	}
	return 0, nil
}

func f3() (int, error) {
	//just returning error when called
	fmt.Println("stacktrace : ", string(debug.Stack()))

	return 0, ErrFunc3
}
