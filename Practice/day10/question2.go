package main

import (
	"fmt"
)

// q3. Create an Interface with one method square(int) int
//     Create a type that implements this interface
//     Create a function Operation that can call the square method using the interface

//     In the main function, create a nil pointer to the concrete type
//     Pass the value to the operation function

//     Operation function calls the method that implements the interface

//     Try to do recovery from panic at different levels

type Shape interface {
	square(int) int
}

type Val struct {
	val int
}

func (v Val) square(i int) int {
	//doesn't matter here since panic happens before this call
	//defer recoverPanic()
	v.val = i * i
	return v.val
}

func opFunc(inter Shape)  {

	fmt.Println("Calling the square method with val", inter)
	//this helps main execute until end of main
	defer recoverPanic()
	inter.square(1)
	fmt.Println("end of operational func")

}

func recoverPanic() {
	msg := recover()
	if msg != nil {
		fmt.Println("Panic Happened")
		fmt.Println(msg)
	}
}

func main() {
	// var in Shape
	// in = nil
	//if defer here it will not execute fmt.Println("End of main")
	//defer recoverPanic()
	//var in Shape
	//in = nil

	var st *Val
	fmt.Println("IN main function")

	opFunc(st)
	fmt.Println("End of main")
}
