package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
)

func (h *Handler) Webhook(c *gin.Context) {
	traceId := uuid.NewString()
	fmt.Println(traceId)

	const MaxBodyBytes = int64(65536)

	// Limit the request body size
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxBodyBytes)

	var event stripe.Event
	err := c.ShouldBindJSON(&event)
	if err != nil {
		slog.Error("Failed to bind JSON", slog.Any("error", err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	fmt.Println(event.Type, "********")
}
