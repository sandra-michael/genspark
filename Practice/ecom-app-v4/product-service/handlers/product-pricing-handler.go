package handlers

import (
	"log/slog"
	"net/http"
	"product-service/internal/products"
	"product-service/pkg/ctxmanage"
	"product-service/pkg/logkey"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getProductOrderDetail(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	productID := c.Param("productID")

	// Extract the context from the HTTP request to pass it to the service layer.
	ctx := c.Request.Context()

	prodOrder, err := h.p.GetStripeProductDetail(ctx, productID)

	if err != nil {
		slog.Error(
			"failed to get stripe price id",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get stripe price id"})
		return
	}
	slog.Info("successfully got stripe customer id", slog.String(logkey.TraceID, traceId))
	c.JSON(http.StatusOK, prodOrder)

}

func (h *Handler) getProductOrderDetails(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	// Extract the context from the HTTP request to pass it to the service layer.
	ctx := c.Request.Context()

	// Parse the request body to get the list of product IDs
	//var productIDs []string

	var req products.ProductOrdersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error(
			"invalid request body",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	// Validate that the list is not empty
	if len(req.ProductIDs) == 0 {
		slog.Error(
			"empty product ID list",
			slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Product ID list cannot be empty"})
		return
	}

	prodOrders, err := h.p.GetStripeProductDetails(ctx, req.ProductIDs)

	if err != nil {
		slog.Error(
			"failed to get stripe price ids",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to get stripe price ids"})
		return
	}
	slog.Info("successfully got stripe customer id", slog.String(logkey.TraceID, traceId))
	c.JSON(http.StatusOK, prodOrders)

}
