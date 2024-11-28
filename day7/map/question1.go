package main

import "fmt"

//q1. Creating and Using Structs in Go
/*
   Todo:
    - Define a struct 'Book' with fields 'Title', 'Author', and 'Pages'.

    - Create a function 'Read' which takes a pointer to a Book, and a number of pages to read.
     After reading, the function should subtract the read pages from the total pages of the book.
      Account for the case where the number of pages to read is more than the pages in the book.
      In that case, just set the book's pages to 0.

    - In the main function, create a new book and initialize it with a title, author, and the number of pages.
    - Then, read some pages from the book and print the title, author, and remaining pages.

   Hint: Use a pointer receiver in the 'Read' method to reflect changes on the original struct.
*/

type Books struct {
	title  string
	author string
	pages  int
}

func (book *Books) reading(readPage int) {
	markRead := book.pages - readPage

	if markRead < 0 {
		book.pages = 0
		return
	}
	book.pages = markRead
}

func (b Books) printBook() {
	fmt.Println(b)

}

func main() {
	awesomeBook := Books{"the best title", "the best author", 120}

	awesomeBook.printBook()

	awesomeBook.reading(3)
	awesomeBook.printBook()

	//completely read and chec use case
	//the case where the number of pages to read is more than the pages in the book.
	//In that case, just set the book's pages to 0.
	awesomeBook.reading(300)
	awesomeBook.printBook()

}
