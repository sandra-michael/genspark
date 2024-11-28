// {
// 	"name": "Alice",
// 	"age": 25
//   }

//   C:\Users\Alice\Documents\example.txt

//   Store above data in string

package main

import (
	"fmt"
)

func main() {

	jsonVal := "{\n\"name\":\"Alice\",\n\"age\":25\n}"

	fmt.Println(jsonVal)

	docVal := " C:\\Users\\Alice\\Documents\\example.txt"
	fmt.Println(docVal)

	//Using Rawstring using backtick

	jsonRaw := `{
		"name": "Alice",
		"age": 25
	}`

	fmt.Println(jsonRaw)

}
