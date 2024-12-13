package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type ContextKey string

const ReqIdKey ContextKey = "reqId"

func main() {
	http.HandleFunc("/a", Mid(HelloHandler))
	//http.Handle("/jsonreq", receieveJson)

	// Start the HTTP server on port 8086 and handle incoming requests
	err := http.ListenAndServe(":8085", nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done")
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Just saying hello from middle wear"))
}

func Mid(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t := time.Now()
		id := uuid.NewString()
		ctx := r.Context()                         // fetching the ctx object from the request
		ctx = context.WithValue(ctx, ReqIdKey, id) // creating an updated ctx with a traceId store in it
		r = r.WithContext(ctx)
		f(w, r)                       //handler function would be called here that was passed to mid
		val := ctx.Value(ReqIdKey)    // fetch the request id for the key
		reqId, ok := val.(ContextKey) // checking if the values exist and of correct type
		if !ok {
			reqId = "unknown"
		}
		log.Printf("reqId, method, url: %s, %s, %s\n", reqId, r.Method, r.URL)
		defer log.Printf("reqId,  duration: %s, %s\n", reqId, time.Since(t))
	}
}
