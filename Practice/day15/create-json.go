package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// struct fields must be exported, so json package can work on it
type person struct {
	FirstName    string          `json:"first_name"`    // json is a field level tag, used by the json package
	Password     string          `json:"-"`             // - means omit the field from the json
	PasswordHash string          `json:"password_hash"` // setting name of the field in the json output
	Perms        map[string]bool `json:"perms"`
}

func main() {
	p := []person{
		{FirstName: "John", Password: "<abc>", PasswordHash: "<#$%^&*(*&^%$>", Perms: map[string]bool{"admin": true}},
		{FirstName: "Jane", Password: "<xyz>", PasswordHash: "<$%^&*()>", Perms: map[string]bool{"admin": false}},
		{FirstName: "Bob", Password: "<qwerty>", PasswordHash: "<$%^&*(>", Perms: map[string]bool{"admin": true}},
	}

	//jsonData, err := json.Marshal(p)
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		log.Println(err)
		return
	}

	// 4 -> Read
	// 2- > write
	// 1 -> execute
	fmt.Println(string(jsonData))
	err = os.WriteFile("data.json", jsonData, 0644)
	if err != nil {
		log.Println(err)
		return
	}

	//f, err := os.OpenFile("data.json", os.O_CREATE|os.O_WRONLY, 0777)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//defer f.Close()
	//n, err := f.Write(jsonData)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}
	//fmt.Println(n, "bytes written")

}
