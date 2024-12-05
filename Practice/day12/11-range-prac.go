package main

import (
	"fmt"
	"sync"
)

func main() {
	wg := new(sync.WaitGroup)
	ch := make(chan int)
	wg.Add(1)
	go func() {
		defer wg.Done()
		//me making sure the inner go routine executes in sync before close
		wgFan := new(sync.WaitGroup)
		//making it work completely async

		for i := 1; i <= 5; i++ {
			wgFan.Add(1)
			//wg.Add(1)
			// fan out pattern, spinning up n number of goroutines, for n number of task
			go func() {
				defer wgFan.Done()
				//defer wg.Done()
				ch <- i
			}()

		}
		wgFan.Wait()
		// close signal range that no more values be sent and it can stop after receiving remaining values
		// close the channel when sending is finished

		// we can't send more values after a channel is closed
		//ch <- 6 // panic: send on closed channel // channel is closed
		close(ch)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}()

	wg.Wait()
}
