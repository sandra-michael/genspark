package handlers

import (
	"log/slog"
	"net/http"
	"product-service/pkg/ctxmanage"
	"product-service/pkg/logkey"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getProductOrderDetail(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	productID := c.Param("productID")

	// Extract the context from the HTTP request to pass it to the service layer.
	ctx := c.Request.Context()

	prodOrder, err := h.p.GetStripeProductDetails(ctx, productID)

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
