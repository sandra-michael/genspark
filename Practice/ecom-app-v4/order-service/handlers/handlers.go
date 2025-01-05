package handlers

import (
	"order-service/gen/proto"
	"order-service/internal/auth"
	"order-service/internal/orders"
	"order-service/internal/stores/kafka"
	"order-service/middleware"
	"os"

	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
)

type Handler struct {
	client      *consulapi.Client
	o           *orders.Conf
	k           *kafka.Conf
	protoclient proto.ProductServiceClient
}

func NewHandler(client *consulapi.Client, o *orders.Conf, k *kafka.Conf, protoclient proto.ProductServiceClient) *Handler {
	return &Handler{client: client, o: o, k: k, protoclient: protoclient}
}

func API(endpointPrefix string, k *auth.Keys, client *consulapi.Client, o *orders.Conf, kafkaConf *kafka.Conf, protoclient proto.ProductServiceClient) *gin.Engine {
	r := gin.New()
	mode := os.Getenv("GIN_MODE")
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	m, err := middleware.NewMid(k)
	if err != nil {
		panic(err)
	}

	h := NewHandler(client, o, kafkaConf, protoclient)
	r.Use(middleware.Logger(), gin.Recovery())

	r.GET("/ping", HealthCheck)
	v1 := r.Group(endpointPrefix)
	{
		v1.POST("/webhook", h.Webhook)
		v1.Use(m.Authentication())
		v1.POST("/checkout/:productID", h.Checkout)
		v1.POST("/checkout/v2/:productID", h.CheckoutWithGrpc)
		v1.POST("/cartcheckout/v2/", h.CartCheckout)
		v1.GET("/ping", HealthCheck)
	}

	return r
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type UserServiceResponse struct {
	StripCustomerId string `json:"stripe_customer_id"`
}

type LineItem struct {
	ProductId string `json:"productId" binding:"required"`
	Quantity  uint64 `json:"quantity" binding:"required"`
}

// Product Order request
type ProductOrdersRequest struct {
	LineItems []LineItem `json:"lineItem" binding:"required"`
}
type ProductServiceResponse struct {
	ProductID string `json:"product_id"`
	Stock     int    `json:"stock"`
	PriceID   string `json:"price_id"`
}

// Product Order request
type CartOrderRequest struct {
	LineItems []LineItem `json:"lineItems" binding:"required"`
}
