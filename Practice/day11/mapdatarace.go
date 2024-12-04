package main

import "fmt"

// q3. Create 4 functions
//     Add(int,int),Sub(int,int),Divide(int,int), CollectResults()
//     Add,Sub,Divide do their operations and push value to an unbuffered channel

//     CollectResult() -> It would receive the values from the channel and print it

var Cal map[string]int = make(map[string]int)

func main() {
	val1 := 1
	val2 := 3
	Add(val1, val2)
	Sub(val1, val2)
	Divide(val1, val2)
	CollectResult()
}

func Add(val1 int, val2 int) {
	res := val1 + val2
	Cal["add"] = res

}

func Sub(val1 int, val2 int) {
	res := val1 - val2
	Cal["sub"] = res

}

func Divide(val1 int, val2 int) {
	res := val1 / val2
	Cal["div"] = res

}

func CollectResult() {
	for key, value := range Cal {
		fmt.Println("Key:", key, "Value:", value)
	}
}
