package handlers

import (
	"fmt"
	"net/http"
	"os"
	"product-service/internal/auth"
	"product-service/internal/products"
	"product-service/middleware"
	"product-service/pkg/ctxmanage"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	consulapi "github.com/hashicorp/consul/api"
)

type Handler struct {
	client   *consulapi.Client
	p        *products.Conf
	validate *validator.Validate
}

func NewHandler(client *consulapi.Client, p *products.Conf) *Handler {
	return &Handler{
		client:   client,
		p:        p,
		validate: validator.New(),
	}
}

func API(client *consulapi.Client, p *products.Conf, k *auth.Keys) *gin.Engine {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	m := middleware.NewMid(k)

	h := NewHandler(client, p)

	prefix := os.Getenv("SERVICE_ENDPOINT_PREFIX")
	if prefix == "" {
		panic("SERVICE_ENDPOINT_PREFIX is not set")
	}

	//DONE create middleware
	r.Use(gin.Logger(), gin.Recovery(), middleware.Logger())

	r.GET("/ping", healthCheck)

	v1 := r.Group(prefix)
	{

		//v1.GET("/stock/:productID"
		//fetch PriceId , stock
		v1.GET("/stock/:productID", h.getProductOrderDetail)
		v1.GET("/", h.fetchAllProducts)
		v1.POST("/stock", h.getProductOrderDetails)

		v1.Use(m.Authentication())

		v1.POST("/", m.Authorize(h.createProduct, auth.RoleAdmin))

		v1.PATCH("/:productID", m.Authorize(h.updateProduct, auth.RoleAdmin))

		//Cart service calls
		v1.POST("/cart/addtocart", h.addToCart)
		v1.POST("/cart/checkout", h.checkout)
		//TODO TEST THIS CODE
		v1.GET("/cart/fetchcart", h.fetchCartDetails)
		v1.DELETE("/cart/:id", h.deleteCartByID)

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
