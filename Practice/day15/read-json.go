package main

import (
	"bufio"
	"encoding/json"
	//"encoding/json"
	"fmt"
	"os"
)

// read the json, written at data.json file
// use json.Unmarshal to convert the byte data to a struct
// os.ReadFile, Scanner, (os.OpenFile -> f.Read)

// struct fields must be exported, so json package can work on it
type person struct {
	FirstName    string          `json:"first_name"`    // json is a field level tag, used by the json package
	Password     string          `json:"-"`             // - means omit the field from the json
	PasswordHash string          `json:"password_hash"` // setting name of the field in the json output
	Perms        map[string]bool `json:"perms"`
}

func printLastNLines(lines []string, num int) []string {
    var printLastNLines []string
    for i := len(lines) - num; i < len(lines); i++ {
        printLastNLines = append(printLastNLines, lines[i])
    }
    return printLastNLines
}

func main(){

	file,err := os.Open("data.json")

	if err != nil {
		fmt.Println("error opening the file ")
		return
	}
	defer file.Close()


	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan(){
		lines = append(lines, scanner.Text())
		
	}

	if err := scanner.Err(); err != nil {
        fmt.Println("Error reading file:", err)
		return 
    }
	//json.Unmarshal("data.json",)

	p, err := json.Unmarshal(file.)


	
    // print the last 10 lines of the file
    // printLastNLines := printLastNLines(lines, 3)
    // for _, line := range printLastNLines {
    //     fmt.Println(line)
    //     fmt.Println("________")
    // }

}