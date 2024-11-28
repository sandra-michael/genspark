package main

import "fmt"

var x = 10

func main() {
	var p *int
	updateValue(p)
	fmt.Println(*p) // 10

}

func updateValue(p1 *int) {
	p1 = &x
	fmt.Println(*p1)
}
