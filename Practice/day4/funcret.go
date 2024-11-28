package main

import (
	"log"
	"strconv"
)

//Create a function that converts a string to float64, if it is successful it returns the value otherwise an error

func main() {
	//input := "sdgreg"
	input := "20.33"
	fltval, err := convStrToFloat64(input)
	if err != nil {
		//log.error - info
		log.Println("There was an error with input : ", input, " ,error message:", err)
		return
	}
	log.Println("The converted float value is ", fltval)
}

func convStrToFloat64(val string) (float64, error) {

	flt, err := strconv.ParseFloat(val, 64)
	if err != nil {
		//log statement here // debug log,info
		return 0, err
	}
	return flt, nil
}
