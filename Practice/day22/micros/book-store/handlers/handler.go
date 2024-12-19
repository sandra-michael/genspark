package handlers

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"


)

func SetupGINRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	return r
}

func HandleBookStore() *gin.Engine {

	r := gin.Default()

	// Step 3: Add OpenTelemetry middleware to Gin
	// This will automatically trace all incoming HTTP requests handled by Gin.
	//needs to be the first middleware which is being used
	r.Use(otelgin.Middleware("book-micro"))

	v1 := r.Group("v1/books")
	{
		v1.POST("/createTable", CreateTable)
		v1.POST("/", Insertfunc)
	}
	return r
}

