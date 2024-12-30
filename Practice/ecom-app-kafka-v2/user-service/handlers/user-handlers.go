package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"user-service/internal/stores/kafka"
	"user-service/internal/users"
	"user-service/pkg/ctxmanage"
	"user-service/pkg/logkey"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Signup(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	//1 KB = 1024 bytes
	// Check if the size of the request body is more than 5KB
	if c.Request.ContentLength > 5*1024 {
		// Log error for payload exceeding size limit
		slog.Error("request body limit breached", slog.String(logkey.TraceID, traceId), slog.Int64("Size Received", c.Request.ContentLength))

		// Return a 400 Bad Request status code along with an error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "payload exceeding size limit"})
		return
	}

	var newUser users.NewUser
	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		// Log error and associate it with a trace id for easy correlation
		slog.Error("json validation error", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))

		// Respond with a 400 Bad Request status code and error message
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": http.StatusText(http.StatusBadRequest)})
		return
	}

	err = h.validate.Struct(newUser)
	if err != nil {
		slog.Error("validation failed", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusBadRequest, gin.H{"error": "please provide values in correct format"})
		return
	}

	ctx := c.Request.Context()
	user, err := h.u.InsertUser(ctx, newUser)
	if err != nil {
		slog.Error("error in creating the user", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User Creation Failed"})
		return
	}
	//topic name - "user-service.account-created"
	//key = userId
	//data=userJson
	data, err := json.Marshal(user)
	if err != nil {
		slog.Error("error in sending user to kafka", slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()))
		c.JSON(http.StatusInternalServerError, gin.H{"message": "User Sending to Kafka Failed"})
		return
	}

	h.kafkaConf.ProduceMessage(ctx, kafka.USER_CREATE_TOPIC, []byte(user.ID), data)

	c.JSON(http.StatusOK, user)

}
