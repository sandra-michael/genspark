package handlers

import (
	"order-service/internal/auth"
	"order-service/middleware"
	"os"

	"github.com/gin-gonic/gin"
	consulapi "github.com/hashicorp/consul/api"
)

type Handler struct {
	client *consulapi.Client
}

func NewHandler(client *consulapi.Client) *Handler {
	return &Handler{client: client}
}

func API(endpointPrefix string, k *auth.Keys, client *consulapi.Client) *gin.Engine {
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

	h := NewHandler(client)
	r.Use(middleware.Logger(), gin.Recovery())

	r.GET("/ping", HealthCheck)
	v1 := r.Group(endpointPrefix)
	{
		v1.POST("/webhook", h.Webhook)
		v1.Use(m.Authentication())
		v1.POST("/checkout/:productID", h.Checkout)
		v1.GET("/ping", HealthCheck)
	}

	return r
}

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
