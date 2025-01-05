package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"product-service/handlers"
	"product-service/internal/auth"
	"product-service/internal/consul"
	"product-service/internal/products"
	"product-service/internal/stores/kafka"
	"product-service/internal/stores/postgres"
	"product-service/protohandler"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "product-service/gen/proto"
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

	// grpcErrors := make(chan error)
	// go func() {
	// 	listener, err := net.Listen("tcp", ":5001")

	// 	//send error to channel

	// 	if err != nil {
	// 		grpcErrors <- err // Send error to the channel
	// 		return
	// 	}

	// 	//NewServer creates a gRPC server which has no service registered
	// 	// creating an instance of the server
	// 	s := grpc.NewServer()

	// 	pb.RegisterProductServiceServer(s, &protohandler.ProtoHandler{})

	// 	//exposing gRPC service to be tested by postman
	// 	reflection.Register(s)

	// 	// Start serving requests
	// 	if err := s.Serve(listener); err != nil {
	// 		grpcErrors <- err // Send error to the channel
	// 	}
	// }()

	// select {
	// case err := <-grpcErrors:
	// 	panic(err)

	// }

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
		//------------------------------------------------------//
		//  Setting up Auth layer
		//------------------------------------------------------//
	*/

	slog.Info("main : Started : Initializing authentication support")
	publicPEM, err := os.ReadFile("pubkey.pem")
	if err != nil {
		return fmt.Errorf("reading auth public key %w", err)
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicPEM)
	if err != nil {
		return fmt.Errorf("parsing auth public key %w", err)
	}

	k, err := auth.NewKeys(publicKey)
	if err != nil {
		return fmt.Errorf("initializing auth %w", err)
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
			//TODO dynamically decrement stock
			//for now we are decrementing for one product
			p.DecrementStock(context.Background(), event.ProductId, 1)
			fmt.Println("successfully decremented the stock of the product")

		}
	}()

	/*
		/*
			//------------------------------------------------------//
			//   Setting up GRPC
			//------------------------------------------------------//
	*/

	grpcErrors := make(chan error)
	go func() {
		listener, err := net.Listen("tcp", ":5001")

		//send error to channel

		if err != nil {
			grpcErrors <- err // Send error to the channel
			return
		}

		//NewServer creates a gRPC server which has no service registered
		// creating an instance of the server
		s := grpc.NewServer()

		pb.RegisterProductServiceServer(s, protohandler.NewProtoHandler(p))

		//exposing gRPC service to be tested by postman
		reflection.Register(s)

		// Start serving requests
		if err := s.Serve(listener); err != nil {
			grpcErrors <- err // Send error to the channel
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
		Handler: handlers.API(p, k),
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
	case err := <-grpcErrors:
		return fmt.Errorf("Grpc error %w", err)
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
