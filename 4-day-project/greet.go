package main

import (
	"log"
	"os"
	"strconv"
)

func main() {
	greet()
}

func greet() {
	data := os.Args[1:] // 1st index till the end
	if len(data) != 3 {
		log.Println("please provide name age marks")
		//os.Exit(1) // quit the program
		return // stops the exec of the current func
	}
	//var err error // default value is nil // it indicates no error happened
	//if err != nil {
	//	// handle the error here
	//}

	name := data[0]
	ageString := data[1]
	marksString := data[2]

	age, err := strconv.Atoi(ageString)
	if err != nil {
		log.Println("invalid age", err)
		return 
	}
	
	mark,err := strconv.Atoi(marksString)
	if err != nil {
		log.Println("Invalid Marks Value",err)
		return
	}
	log.Println(name, age,mark)
	//strconv.ParseFloat()
}
