package main

import (
	"fmt"
	"sync"
)

// q3. Create a bookCab function that takes the name of the user trying to book a cab
//     Assume only one user can book a cab at a time
//     Create a global variable to hold the number of cabs available

//     Check if a cab is available, if yes print a msg cab is available otherwise unavailable

//     In the main function, run a loop 5 times and call bookCab function as goroutine to simulate
//     multiple users are trying to book a cab

var AvailableCabs = 5
func main(){
	users := []string{"san","jam","roxy","snow","ruby","mike"}
	wg := new(sync.WaitGroup)
	for _ , val := range users {
		wg.Add(1)
		go bookCab(val,wg)
	}
	wg.Wait()

}

func bookCab(userName string, wg *sync.WaitGroup){
	defer wg.Done()
	if AvailableCabs > 0 {
		AvailableCabs = AvailableCabs - 1
		fmt.Println("Cab available for user ",userName)
		return
	}
	fmt.Println("Cab unavailable for user ",userName," sorry try again")

		
}