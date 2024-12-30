package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	api := http.Server{
		Addr:              ":8082",
		ReadHeaderTimeout: time.Second * 200,
		WriteTimeout:      time.Second * 200,
		IdleTimeout:       time.Second * 200,
		Handler:           handle(),
	}

	api.ListenAndServe()
}

func handle() *gin.Engine {

	r := gin.Default()

	r.GET("/*path", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "hi "})
	})
	return r
}
