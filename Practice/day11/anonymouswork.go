package main

import (
	"fmt"
	"sync"
	"time"
)

// q2. Follow up to the previous question
//     Spin up one anonymous goroutine in the work function
//     This goroutine prints some stuff on the screen and sleeps for 100ms
//     Make sure you wait for every goroutine to finish and end everything gracefully

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
	//without defer the Done code wasn't always executing end statement
	defer wg.Done()

	fmt.Println("Work with workId : ", id, " is going on")
	//this makes my wait wait for the anonymous go routine
	//else we can just add two routines in wg.Add(2) and defer wg.Done() inside the go routine
	wgA := new(sync.WaitGroup)

	wgA.Add(1)
	go func() {
		//trying without defer
		wgA.Done()
		//need to verify if defer is needed in the internal statement
		//defer wgA.Done()
		fmt.Println("Hi I'm Mr anonymous working for id", id)
		time.Sleep(100 * time.Millisecond)

	}()
	wgA.Wait()
	fmt.Println("End of Work with workId : ", id)
}
