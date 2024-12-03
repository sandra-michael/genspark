package main

import (
	"errors"
	"fmt"
)

// q1.	create a program that manages a collection of books and number of books available
//     allow users to search for books by title.
// 	The program should handle errors gracefully if a book is not found or if there are any issues accessing the collection.

// Use a map to store book Name and their counter.
// Functionality:
// Implement
//   - AddBook(title string,counter int) error
//     -to add a new book to the collection.
//
// FetchBookCounter(name) (int, error)
//
//	-to retrieve a book by its name.
//
// Error Handling:
// Use a struct to handle error
// User errors.As in main to check if struct is present inside the chain or not
var ErrorNotFound = errors.New("Not FOund")
var ErrorBookExists = errors.New("Book Already Exists")

type FindError struct {
	fname string
	title string
	err   error
}

func (f *FindError) Error() string {
	return "BookStr." + f.fname + ": " + "fetching " + f.title + ": " + f.err.Error()
}

type BookStr struct {
	title   string
	counter int
}

var BookDb map[string]BookStr = make(map[string]BookStr)

func (b BookStr) AddBook(title string, counter int) error {
	_, ok := BookDb[title]
	if !ok {
		//SImple
		//in this case errors.as wont work
		//return 0, ErrorBookExists
		//Dynameic
		return &FindError{"AddBook", title, ErrorBookExists}
	}
	BookDb[title] = BookStr{title, counter}
	return nil
}

func (b BookStr) FetchBookCounter(title string) (int, error) {

	book, ok := BookDb[title]

	if ok {
		//SImple
		//in this case errors.as wont work
		//return 0, ErrorNotFound
		//Dynameic
		return 0, &FindError{"FetchBookCounter", title, ErrorNotFound}
	}
	return book.counter, nil

}

func main() {

	var b1 BookStr

	b1.AddBook("awesome book", 100)

	//uncomment below for failure case
	//val, err := b1.FetchBookCounter("some")
	//uncomment below for success case
	val, err := b1.FetchBookCounter("awesome book")

	if err != nil {
		//Checking using errors,As
		var fe *FindError
		if errors.As(err, &fe) {
			fmt.Println("error function name ", fe.fname)
			fmt.Println("error book title name ", fe.title)
			fmt.Println("Complete stack trace", err)
			return
		}
		fmt.Println("error : ", err)
		return
	}
	fmt.Println("Book found ", val)

	//trying to add a book which exists
	err = b1.AddBook("awesome book", 100)
	if err != nil {
		//Checking using errors,As
		var fe *FindError
		if errors.As(err, &fe) {
			fmt.Println("error function name ", fe.fname)
			fmt.Println("error book title name ", fe.title)
			fmt.Println("Complete stack trace", err)
			return
		}
		fmt.Println("error : ", err)
		return
	}

}
