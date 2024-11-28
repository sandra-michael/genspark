package main

import "fmt"

// Q2. Create an initial list of users.
//     Create a new slice named as emp from a part users slice.
//     Check both slices length and capacity.

//     Append some more user names to the *user* slice.
//     Check len, cap again for both slices
//     Print the values as well

//     Append some emp names to the *emp* slice
//     Check len, cap again for both slices
//     Print the values as well

//     Compare the contents and capacities of the original list and the modified slice.
//     Try to understand and visualize what is happening

func checkSlice(name string, s []string) {
	fmt.Println("\nSlice name: ", name, " slice: ", s)
	fmt.Println("Slice name: ", name, " len: ", len(s), " cap: ", cap(s), " memory: ", &s[0])
}

func main() {

	user := []string{"a", "b", "c", "d", "e", "f", "g"}

	emp := user[3:6]

	checkSlice("user", user)
	checkSlice("emp", emp)

	fmt.Println("\n\n--------Appending names to user ---------------")

	user = append(user, "l", "m", "n")

	checkSlice("user", user)
	checkSlice("emp", emp)

	fmt.Println("\n\n-----------Append some emp names to the *emp* slice----------")

	emp = append(emp, "o", "p")

	checkSlice("user", user)
	checkSlice("emp", emp)
}
