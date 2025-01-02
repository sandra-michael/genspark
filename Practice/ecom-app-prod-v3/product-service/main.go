package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"product-service/handlers"
	"product-service/internal/consul"
	"product-service/internal/products"
	"product-service/internal/stores/kafka"
	"product-service/internal/stores/postgres"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	setupSlog()
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("error in loading env file")
	}

	err = startApp()
	if err != nil {
		panic(err)
	}
}

func startApp() error {
	slog.Info("Migrating tables for user-service if not already done")
	db, err := postgres.OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	err = postgres.RunMigrations(db)
	if err != nil {
		return err
	}

	p, err := products.NewConf(db)
	if err != nil {
		return err
	}

	/*
		/*
			//------------------------------------------------------//
			//   Consuming Kafka TOPICS [ORDER SERVICE EVENTS]
			//------------------------------------------------------//
	*/
	go func() {
		ch := make(chan kafka.ConsumeResult)
		go kafka.ConsumeMessage(context.Background(), kafka.TopicOrderPaid, kafka.ConsumerGroup, ch)
		for v := range ch {
			if v.Err != nil {
				fmt.Println(v.Err)
				continue
			}
			fmt.Printf("Consumed message: %s", string(v.Record.Value))
			var event kafka.OrderPaidEvent
			json.Unmarshal(v.Record.Value, &event)
			// create a method over internal/products to decrement the stock value by quantity
			fmt.Println("decrement the stock of the product")
			p.UpdateProducts(context.Background(), event.ProductId)
			fmt.Println("successfully decremented the stock of the product")

		}
	}()

	//setting up http server
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
		//port = "8083"
	}

	api := http.Server{
		Addr:         ":" + port,
		ReadTimeout:  8000 * time.Second,
		WriteTimeout: 800 * time.Second,
		IdleTimeout:  800 * time.Second,
		//handlers.API returns gin.Engine which implements Handler Interface
		Handler: handlers.API(p),
	}
	serverErrors := make(chan error)
	go func() {
		serverErrors <- api.ListenAndServe()
	}()

	/*
			//------------------------------------------------------//
		               Registering with Consul
			//------------------------------------------------------//
	*/

	consulClient, regId, err := consul.RegisterWithConsul()
	if err != nil {
		return err
	}

	defer consulClient.Agent().ServiceDeregister(regId)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM, os.Kill)
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error %w", err)
	case <-shutdown:

		fmt.Println("Shutting down server gracefully")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		//Shutdown gracefully shuts down the server without interrupting any active connections.
		//Shutdown works by first closing all open listeners, then closing all idle connections,
		err = api.Shutdown(ctx)
		if err != nil {

			//forceful closure
			err := api.Close()
			if err != nil {
				// returning error to main if everything fails, the main would panic
				return fmt.Errorf("could not stop server gracefully %w", err)
			}
		}
	}

	return nil

}

func setupSlog() {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		//AddSource: true: This will cause the source file and line number of the log message to be included in the output
		AddSource: true,
	})

	logger := slog.New(logHandler)
	//SetDefault makes l the default Logger. in our case we would be doing structured logging
	slog.SetDefault(logger)
}
