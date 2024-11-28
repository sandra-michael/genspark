package slice

import "fmt"

// q1. Write a Go program that:
//     Creates an empty slice with an initial capacity of 1.
//     Appends integers from 1 to  to the slice.
//     Tracks and prints the capacity change every time the slice's capacity increases.
//     Prints the total number of capacity changes at the end.

//     Formula:= (currentCap-lastCap) / lastCap * 100
//     // Hint :- use type casting

func SliceCap(low int, high int) int {
	var s []int
	capChange := 0
	lastCap := cap(s)

	for i := low; i < high; i++ {
		s = append(s, i)
		currentCap := cap(s)

		if lastCap < currentCap {
			var percentage = float64(currentCap-lastCap) / float64(lastCap) * 100
			fmt.Println("Capacity Changed from ", lastCap, " to ", currentCap, "Cap increase persentage", percentage)
			lastCap = currentCap
			capChange = capChange + 1
		}

	}

	//lastCap := cap(s)
	//Formula:= (currentCap-lastCap) / lastCap * 100
	return capChange
}
