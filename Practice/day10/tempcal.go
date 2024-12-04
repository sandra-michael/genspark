package main

import "fmt"

// q3. Create a new custom type based on float64 for handling temperatures in Celsius.
//     Implement the Following Methods (not functions):
//     Method 1: ToFahrenheit
//     Description: Converts the Celsius temperature to Fahrenheit.
//     Signature: ToFahrenheit() float64
//     Method 2: IsFreezing
//     Description: Checks if the temperature is at or below the freezing point (0°C).
//     Signature: IsFreezing() bool

type Celsius float64

func (c Celsius) ToFahrenheit() float64 {
	//(0°C × 9/5) + 32 = 32°F
	return float64(c*9/5) + 32
}

func (c Celsius) IsFreezing() bool {
	if c <= 0 {
		return true
	}
	return false
}

func main() {
	//var temp Celsius = 10

	//var temp Celsius = -1

	if temp.IsFreezing() {
		fmt.Println("It's cold ")
	}

	if !temp.IsFreezing() {
		fmt.Println("It's HOT!!!!!!!!!!!!!!!!!!!!!!!")
	}
	fmt.Println("Converted to farhenhite", temp.ToFahrenheit())

}
