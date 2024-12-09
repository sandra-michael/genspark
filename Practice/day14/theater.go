package main

import (
	"fmt"
	"sync"
	"time"
)

// q1.  Make a struct Theater with the following fields: Seats(int=1), RWMutex, userName chan string.

//      Create two methods over a struct

//      The first method book a seat in the theater. If the seat is equal to zero, no one can book it.
//      ( In the booking method, put simple print statements that show booking has been made if seats are available)

//      Once the seat is booked in Theater, add the name of the user in the userName channel.
//      Create a second Method, printInvoice(),  It fetches the userName from the channel and print it.

//     Note:-
//      Call the bookSeats & printInvoice() method as a goroutine in the main function.
//      For example:-

//      for i:=1; i<=3; i++ {
//           go t.bookSeats()
//      }
//      go t.printInvoice()

//      The program should quit gracefully without deadlock.

type theater struct {
	seat int
	//adding mutex to struct and it need not be a pointer if methods have pointers 
	
}

func (t *theater) bookSeats(c chan<- string, name string, m *sync.RWMutex, wg *sync.WaitGroup) {
	m.Lock() //locking my resource so it is not available for any other go routine
	defer m.Unlock()
	defer wg.Done()
	//fmt.Println("Available seats", t.seat)

	if t.seat < 1 {
		fmt.Println("seat is not available for", name)
	} else {
		fmt.Println("seat is available for", name)
		time.Sleep(5 * time.Second)
		fmt.Println("booking confirmed", name)
		t.seat--
		c <- name
	}

}

//trying without a mutex

func (t *theater) bookSeatsWithoutMutex(c chan<- string, name string, wg *sync.WaitGroup) {
	//m.Lock() //locking my resource so it is not available for any other go routine
	//defer m.Unlock()
	defer wg.Done()
	if t.seat < 1 {
		fmt.Println("seat is not available for", name)
	} else {

		fmt.Println("seat is available for", name)
		time.Sleep(5 * time.Second)
		fmt.Println("booking confirmed", name)
		t.seat--
		c <- name
	}

}
func (t *theater) checkSeats(m *sync.RWMutex, wg *sync.WaitGroup, i int) {
	m.RLock()
	defer m.RUnlock()
	defer wg.Done()
	fmt.Println("Available seats", t.seat, " call ", i)
}

func (t *theater) printInvoice(c <-chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range c {
		fmt.Println("invoice", i)
	}
}

func main() {
	//m := new(sync.Mutex)
	m := new(sync.RWMutex)
	wg := new(sync.WaitGroup)

	//creating a new channel so we can close once the sender sends all users
	wgUserChan := new(sync.WaitGroup)

	var t theater
	t.seat = 2
	userName := make(chan string, 0)
	user := []string{"san", "jam", "roxy"}

	for _, val := range user {
		time.Sleep(5 * time.Second)
		wgUserChan.Add(1)
		go t.bookSeats(userName, val, m, wgUserChan)
		//go t.bookSeatsWithoutMutex(userName, val, wg)

	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		wgUserChan.Wait()
		close(userName)
	}()

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go t.checkSeats(m, wg, i)
		//go t.bookSeatsWithoutMutex(userName, val, wg)
	}
	wg.Add(1)
	go t.printInvoice(userName, wg)

	wg.Wait()
}
