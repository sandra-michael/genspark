package main

import (
	"fmt"
	"sync"
	//"time"
)

// q1. create a function work that takes a work id and print work {id} is going on
//     In the main function run a loop to run work function 10 times
//     make the work function call concurrent
//     Make sure your program waits for work function to finish gracefully

func main() {
	fmt.Println("Start of main")
	wg := new(sync.WaitGroup)

	wg.Add(10)
	for i := 1; i < 11; i++ {
		go work(i, wg)
	}

	wg.Wait()
	//trying to execute all go routines by putting main to sleep
	//this works executing all 10
	//time.Sleep(time.Second)

	fmt.Println("End of main")
}

func work(id int, wg *sync.WaitGroup) {
	wg.Done()

	fmt.Println("Work with workId : ", id, " is going on")
	fmt.Println("End of Work with workId : ", id)
}
