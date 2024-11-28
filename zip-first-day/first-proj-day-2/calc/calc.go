package calc

//good practive package name is same as folder name
import "fmt"

func Add(a, b) { // if firsrt letter of func is in Uppercase the func would be exported
	fmt.Println(a + b)

}
