package main

import (
	"fmt"
	"strings"
)

// q1. Create a struct (Author)
//     Two Field:- Name, Books[slice]
//     Create two methods, one appends new books to the book slice , other prints the struct

//     Create a function that accepts the struct and append values to the book slice

//     Create a function that would accept the Books field, not the struct and append some more books

type Author struct {
	name  string
	books []string
}

// Creating methods meaning these can only be called by slice
func (a *Author) appendNewBooks(newBooks ...string) {
	a.books = append(a.books, newBooks...)
}

func (a *Author) printAuthorStrut() {
	//fmt.Printf("%+v \n ", a)
	//fmt.Println("getting the pointer to strct", a)
	//output of above getting the pointer to strct &{somoneAwesome [a book another book yet another again a book 2 a book 3]}
	fmt.Println("Author name : ", a.name, " list of books:  ", strings.Join(a.books, " , "))
}

func main() {
	//a1 := Author{"somoneAwesome", []string{"a book","another book"} }
	var a1 Author
	a1.name = "somoneAwesome"
	a1.books = []string{"a book", "another book"}
	a1.printAuthorStrut()
	//Enen though using pointer does not require &
	a1.appendNewBooks("yet another", "again")
	a1.printAuthorStrut()

	//Since its not a method requires &
	appendNewBooks(&a1.books, "a book 2", "a book 3")

	a1.printAuthorStrut()
}

func appendNewBooks(books *[]string, newBooks ...string) {
	*books = append(*books, newBooks...)
}
