package main

import "fmt"

// Create a function that takes a list of users
// This func can append new values to the list or change the existing elems
// But Make sure this function can't modify the original slice
// that was created in the main function
func slicerModify(s []int) []int {
	s[0] = 111
	s[1] = 222
	s = append(s, 11, 12, 13, 14)
	return s
}

func checkSlice(name string, s []int) {
	fmt.Println("\nSlice name: ", name, " slice: ", s)
	fmt.Println("Slice name: ", name, " len: ", len(s), " cap: ", cap(s), " memory: ", &s[0])
}

func main() {

	og := []int{1, 2, 3, 4, 5, 6, 7}

	//creating a deep copy of og

	dc := make([]int, len(og), cap(og))

	fmt.Println("\n\n-------checking original  ------------")
	checkSlice("og", og)
	

	copy(dc, og)

	fmt.Println("\n\n-------checking deep copy  ------------")
	checkSlice("dc", dc)
	
	dc = slicerModify(dc)

	fmt.Println("\n\n-------checking original after change in copy ------------")
	checkSlice("og", og)
	fmt.Println("\n\n-------checking deep copy after changey ------------")
	checkSlice("dc", dc)
}
