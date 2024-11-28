package main

import "fmt"

//Declare a variable to represent temperature in Celsius. Convert this temperature to Fahrenheit using the formula
//try to type cast

func main() {
	var celcius = 32
	var fahrenheit float64 = float64((celcius * 9 / 5) + 32)
	fmt.Printf("%d celcius is %f Fahrenheit", celcius, fahrenheit)

}
