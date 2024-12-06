package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

// q1. Create a slice with 3 random urls
//     Create a function doGetRequest()
//     doGetRequest:
//         It spins up 2 goroutines
//         1st goroutines do get request and put the url,resp,err inside one single channel
//         //1st goroutine spins up n number of goroutines for n number of urls (fanout pattern)
//         2nd goroutines fetch the values from the channel and perform following operations
//             -check err
//             -read body
//             -check if status code above 299
//             -and print url resp.Status

var urls []string = []string{"https://pkg.go.dev/", "https://pkg.go.dev/search?q=string", "https://pkg.go.dev/strconv"}

type response struct {
	url  string
	resp *http.Response
	err  error
}

func main() {
	wg := new(sync.WaitGroup)

	doGetRequest(wg)
	fmt.Println("end of main")
}

func doGetRequest(wg *sync.WaitGroup) {

	//creating a worker group
	wgWorker := new(sync.WaitGroup)
	getch := make(chan response)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for _, val := range urls {
			wgWorker.Add(1)
			go func() {
				fmt.Println("sending response for url", val)
				defer wgWorker.Done()
				res, err := http.Get(val)
				str := response{url: val, resp: res, err: err}
				getch <- str
			}()

		}
		wgWorker.Wait()
		close(getch)

	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range getch {
			fmt.Println("Url ", v.url)
			defer v.resp.Body.Close()
			_, err := io.ReadAll(v.resp.Body)
			if err != nil {
				log.Println(err)
				//this quits the go routines hence the rest will quit as well
				//return
				continue
			}
			//resStr := string(bytes)
			//fmt.Println("res ", resStr[0:30])
			if v.resp.StatusCode > 299 {
				fmt.Println("res status code is above 299", v.resp.StatusCode)
				continue
			}
			fmt.Println("res status", v.resp.Status)
		}
	}()
	wg.Wait()

	fmt.Println("do Get req end")

}
