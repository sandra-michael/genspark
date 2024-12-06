package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	wg := new(sync.WaitGroup)

	workerwg := new(sync.WaitGroup)

	// select is used when we want to listen or send values to over a multiple channel
	get := make(chan string)
	post := make(chan string)
	put := make(chan string)

	done := make(chan struct{}) // datatype doesn't matter

	workerwg.Add(1)
	go func() {
		defer workerwg.Done()
		time.Sleep(1 * time.Second)
		get <- "get"

	}()

	workerwg.Add(1)
	go func() {
		defer workerwg.Done()
		time.Sleep(50 * time.Millisecond)
		post <- "post"
	}()

	workerwg.Add(1)
	go func() {
		defer workerwg.Done()
		put <- "put"
		put <- "p1"

	}()

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		workerwg.Wait()
		close(done)
		
	}()

	// not efficient // because we have to wait for get even if it is taking long time execute
	//fmt.Println(<-get)
	//fmt.Println(<-post)
	//fmt.Println(<-put)

	// the problem with below loop is if less or more values are sent it would not work
	// and deadlock would happened
	//for i := 0; i < 3; i++ {
	//	// whichever case is not blocking exec that first
	//	//whichever case is ready first, exec that.
	//	// possible cases are chan recv , send , default
	//	select {
	//	case g := <-get:
	//		fmt.Println(g)
	//	case p := <-post:
	//		fmt.Println(p)
	//	case pu := <-put:
	//		fmt.Println(pu)
	//
	//	}
	//}
	//close(done) // close is a send signal, and select can recv it
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case g := <-get:
				fmt.Println(g)
			case p := <-post:
				fmt.Println(p)
			case pu := <-put:
				fmt.Println(pu)
			case <-done:
				fmt.Println("all values are received")
				return

			}
		}
	}()
	wg.Wait()

}
