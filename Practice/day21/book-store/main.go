package main

import (
	//"book-store/handlers"
	"book-store/handlers"

	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// initialize http service
	//chi, http.DefaultServeMux, gin
	api := http.Server{
		Addr:              ":8082",
		ReadHeaderTimeout: time.Second * 200,
		WriteTimeout:      time.Second * 200,
		IdleTimeout:       time.Second * 200,
		Handler:           handlers.HandleBookStore(),
	}
	// Channel to listen for OS signals (like SIGTERM, SIGINT) for graceful shutdown
	shutdown := make(chan os.Signal, 1)

	// Register the shutdown channel to receive specific system interrupt signals
	signal.Notify(shutdown, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Channel to capture server errors during runtime, like port already being used
	serverError := make(chan error)

	// Goroutine to handle server startup and listen for incoming requests
	go func() {
		serverError <- api.ListenAndServe()
	}()

	// select statement to handle either server errors or shutdown signals
	select {
	// this error would happen if the service is not able to start
	case err := <-serverError:
		// Panic if the server fails to start
		panic(err)
	case <-shutdown:
		fmt.Println("Graceful Shutdown Server...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		//Shutdown gracefully shuts down the server without interrupting any active connections.
		//Shutdown works by first closing all open listeners, then closing all idle connections,
		err := api.Shutdown(ctx)
		if err != nil {
			// force close
			err := api.Close()
			panic(err)
		}

	}

}
