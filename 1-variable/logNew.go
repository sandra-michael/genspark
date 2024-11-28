package main

import (
	"log"
	"os"
)

func main() {
	mylog := log.New(os.Stdout, "my:", log.LstdFlags|log.Lshortfile)
	mylog.Println("from mylog")

}
