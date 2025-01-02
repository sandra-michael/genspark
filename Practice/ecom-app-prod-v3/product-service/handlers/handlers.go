package handlers

import (
	"fmt"
	"net/http"
	"os"
	"product-service/internal/products"
	"product-service/middleware"
	"product-service/pkg/ctxmanage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	p        *products.Conf
	validate *validator.Validate
}

func NewHandler(p *products.Conf) *Handler {
	return &Handler{
		p:        p,
		validate: validator.New(),
	}
}

func API(p *products.Conf) *gin.Engine {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	h := NewHandler(p)

	prefix := os.Getenv("SERVICE_ENDPOINT_PREFIX")
	if prefix == "" {
		panic("SERVICE_ENDPOINT_PREFIX is not set")
	}

	//DONE create middleware
	r.Use(gin.Logger(), gin.Recovery(), middleware.Logger())

	r.GET("/ping", healthCheck)

	v1 := r.Group(prefix)
	{
		v1.POST("/", h.createProduct)

		//v1.GET("/stock/:productID"
		//fetch PriceId , stock
		v1.GET("/stock/:productID", h.getProductOrderDetail)

	}

	//TODO CREATE API
	return r
}

func healthCheck(c *gin.Context) {
	//Need to add context
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	fmt.Println("healthCheck handler ", traceId)
	//JSON serializes the given struct as JSON into the response body. It also sets the Content-Type as "application/json".
	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}
