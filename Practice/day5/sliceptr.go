package main

import "fmt"

// q1. In main function create a slice of names.
//     Add two elements in it.

//     Create a function AddNames which appends new names to the slice
//     Use double pointer concept to make AddNames function work
// NOTE :we know slicer already holds pointers

func main() {
	slc := []string{"name1", "name2"}

	fmt.Println("Address of slice ", &slc[0], " Slice Value before update", slc)

	addNames(&slc) // sending address of slice
	fmt.Println("Address of slice ", &slc[0], " Slice Value after update", slc)

	fmt.Println(slc)
}

func addNames(ipSlc *[]string) {
	fmt.Println("value of pointer ipSlc ", ipSlc)
	fmt.Println("value of pointer deref ipSlc ", *ipSlc)
	fmt.Println("address of ipSlc ", &ipSlc)

	*ipSlc = append(*ipSlc, "name3", "name4")

}
