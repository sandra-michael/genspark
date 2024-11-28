package main

import (
	//"cli/calc"
	"log"
	"os"
	"strconv"

	"cli.go/calc"
)

func main() {
	log.Println("Within main")
	calculator()

}

func calculator() {

	data := os.Args[1:]

	operator := data[0]
	if len(data) != 3 {
		log.Println("please enter 3 values")
		return
	}

	// if !strings.Contains("+-*%", operator) || len(operator) > 1 {
	// 	log.Println("Please enter the correct operator")
	// 	return
	// }

	//Check the values are integers
	firstValueStr := data[1]

	firstVal, err := strconv.Atoi(firstValueStr)
	if err != nil {
		log.Println("please enter a proper value for value one ")
		return
	}

	secondValueStr := data[2]

	secondVal, err := strconv.Atoi(secondValueStr)
	if err != nil {
		log.Println("please enter a proper value for value two ")
		return
	}

	//log.Println(secondVal, firstVal, operator)

	switch operator {
	case "+":
		log.Println(calc.Add(firstVal, secondVal))
	case "-":
		log.Println(calc.Sub(firstVal, secondVal))
	case "*":
		log.Println(calc.Mul(firstVal, secondVal))
	case "%":
		log.Println(calc.Mod(firstVal, secondVal))
	default:
		log.Println("Please enter the correct operator")

	}

}
