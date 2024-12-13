package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"context"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

var users = []string{"sandra", "mike", "diwakar", "alka"}
var posts = []string{"some post", "then post ", "sdvnmlomlolm", "random posts "}

type ContextKey string

const ReqIdKey ContextKey = "reqId"

// add traceId middleware to trace the request to user routes
func main() {
	fmt.Println("in main")
	mux := chi.NewRouter()

	// set global middleware
	mux.Use(middleware.Logger)
	// ctx := mux.Context()                       // fetching the ctx object from the request
	// ctx = context.WithValue(ctx, ReqIdKey, id) // creating an updated ctx with a traceId store in it
	// r = r.WithContext(ctx)

	// localhost:8080/v1/users/123
	mux.Route("/v1/users", func(r chi.Router) {
		r.Use(MidTraceId)
		// get user
		r.Get("/", getUsers)

		//get user by id
		r.Get("/{id}", getUser)

		// create one user
		r.Post("/create", createUser)
	})

	// localhost:8080/v1/posts/123
	mux.Route("/v1/posts", func(r chi.Router) {
		r.Use(middleware.Logger, middleware.Recoverer)
		// fetch all posts
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		// fetch post by id
		r.Get("/{id}", func(w http.ResponseWriter, r *http.Request) {})
		// create post
		r.Post("/create", func(w http.ResponseWriter, r *http.Request) {})
	})

	// Start the HTTP server on port 8086 and handle incoming requests
	err := http.ListenAndServe(":8085", mux)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("DOne")
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	//this will recover and not exit
	// panic("sdffvtfdtehbg Panic")
	//this will exit any manually run go routine panic will exit
	//go panic("sdffvtfdtehbg Panic")
	w.WriteHeader(http.StatusOK)

	data, err := json.Marshal(users)
	if err != nil {
		log.Println("could not marshal", err)
		http.Error(w, "user not found in db", http.StatusNotFound)
		return
	}
	w.Write(data)
	return

}
func getUser(w http.ResponseWriter, r *http.Request) {
	//log.Fatal("Fatal Error ")
	idVal := chi.URLParam(r, "id")
	idCon, err := strconv.Atoi(idVal)
	if err != nil {
		http.Error(w, "Invalid id provide proper integer", http.StatusNotFound)
		return
	}
	if len(users) < idCon {
		http.Error(w, "User Does not exist ", http.StatusNotFound)
		return
	}

	w.Write([]byte("name, " + users[idCon] + "for id" + idVal))
	return

}
func createUser(w http.ResponseWriter, r *http.Request) {
	Name := r.URL.Query().Get("name")

	users = append(users, Name)
	w.Write([]byte("User Created, " + Name))
}

// func mid(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { next.ServeHTTP(w, r) })
// }

// func MidTraceId(next http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		fn := func(w http.ResponseWriter, r *http.Request) {

// 			id := uuid.NewString()
// 			ctx := r.Context()                         // fetching the ctx object from the request
// 			ctx = context.WithValue(ctx, ReqIdKey, id) // creating an updated ctx with a traceId store in it
// 			r = r.WithContext(ctx)                     // putting context inside the request object
// 			next.ServeHTTP(w, r)
// 		}

// 		return http.HandlerFunc(fn) // calling next thing in the chain

// 	}
// }

func MidTraceId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.NewString()
		ctx := r.Context()                         // fetching the ctx object from the request
		ctx = context.WithValue(ctx, ReqIdKey, id) // creating an updated ctx with a traceId store in it
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
