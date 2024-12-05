package main

import (
	"fmt"
	"sync"
)

// q1. Create 4 functions
//     Add(int,int),Sub(int,int),Divide(int,int), CollectResults()
//     Add,Sub,Divide do their operations and push value to an unbuffered channel

//     CollectResult() -> It would receive the values from the channel and print it

//var Cal map[string]int = make(map[string]int)

var ubch chan int = make(chan int)

func main() {
	val1 := 1
	val2 := 3
	wg := new(sync.WaitGroup)

	//spining go routines
	wg.Add(3)
	go Add(val1, val2, wg)
	go Sub(val1, val2, wg)
	go Divide(val1, val2, wg)

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			CollectResult(i)
		}()
	}

	wg.Wait()
	fmt.Println("End of main")
}

func Add(val1 int, val2 int, wg *sync.WaitGroup) {
	defer wg.Done()
	res := val1 + val2
	fmt.Println("Adding values ", val1, " ", val2, " res ", res)
	ubch <- res

}

func Sub(val1 int, val2 int, wg *sync.WaitGroup) {
	defer wg.Done()
	res := val1 - val2
	fmt.Println("sub values ", val1, " ", val2, " res ", res)

	ubch <- res

}

func Divide(val1 int, val2 int, wg *sync.WaitGroup) {
	defer wg.Done()
	res := val1 / val2
	fmt.Println("div values ", val1, " ", val2, " res ", res)

	ubch <- res

}

func CollectResult(i int) {
	//defer wg.Done()
	// for key, value := range Cal {
	// 	fmt.Println("Key:", key, "Value:", value)
	// }
	fmt.Println("Receiver ", i)
	fmt.Println("Rec ", i, " val ", <-ubch)
}
