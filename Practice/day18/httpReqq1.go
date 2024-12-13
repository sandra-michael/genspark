package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// q1. Create a handler which accepts any kind of json and prints it
//     check if client is still connected or not
//     if client is connected then return json processed otherwise just move on

//     Hint: use map[string]any if not sure about json structure
//     The any type has a tradeoff as well, you can't access the fields directly,
//     type assertion needs to be done every time

//     TO find the json what user have sent, you can use r.Body

func main() {
	http.HandleFunc("/jsonreq", receieveJson)
	//http.Handle("/jsonreq", receieveJson)

	// Start the HTTP server on port 8086 and handle incoming requests
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}

func receieveJson(w http.ResponseWriter, r *http.Request) {

	// Extract the "user_name" query parameter from the request URL
	//readJ := strings.NewReader(r.Body())
	//use umarshal here
	reqBody, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("Error reading request")
		return
	}

	var jsonMap map[string]any = make(map[string]any, 0)
	//json.NewDecoder(copyBody).Decode(&jsonMap)
	err = json.Unmarshal(reqBody, &jsonMap)
	if err != nil {
		fmt.Println("Error converting to map", err)
		return
	}

	fmt.Println("re ", string(reqBody))
	fmt.Println("Json input ", jsonMap)

	w.WriteHeader(http.StatusOK)

	for key, val := range jsonMap {
		fmt.Println("key ", key, " val ", val)

		//w.Write([]byte(key + " : " + val.(string)))
	}

	w.Write([]byte("success"))

	//http.Error(w, "not working", http.StatusNotFound)

	return

}
