package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
)

func main() {
	// Register the service with Consul for service discovery.
	registerServiceConsul()

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "Hello microservice health check!")
	})

	r.GET("/hye/:name", func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.String(http.StatusBadRequest, "Name is required")
			return
		}
		c.String(200, "Hello"+" "+name+"!")
	})
	log.Println("Listening on port 80")
	r.Run(":80")
}

// registerServiceConsul registers the HTTP service with Consul for service discovery.
func registerServiceConsul() {
	// Create a default configuration for connecting to Consul.
	config := api.DefaultConfig()

	// Set the address of the Consul server. Change this to point to your actual Consul service.
	config.Address = "http://consul.sandra:8500"

	// Create a new client to interact with Consul.
	consul, err := api.NewClient(config)
	if err != nil {
		// If an error occurs while creating the client, stop the application by panicking.
		panic(err)
	}

	// Fetch the hostname of the machine/service from the environment variable `HOSTNAME`.
	address := os.Getenv("HOSTNAME")
	if address == "" {
		// If `HOSTNAME` is not set, stop the application with an error message.
		panic("HOSTNAME not set")
	}

	// Create a new Consul service registration object.
	registration := &api.AgentServiceRegistration{}

	// Assign a name for the service. This is how other applications will reference this service in Consul.
	registration.Name = "hello-service"

	// Assign a unique ID for the service, combining the service name and the hostname to make it unique.
	registration.ID = "hello-service" + address

	// Assign the hostname (service's address) and port on which this service is running.
	registration.Address = address
	registration.Port = 80

	// Register the service with Consul. If this process fails, stop the application.
	err = consul.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
}
