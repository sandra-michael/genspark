package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"product-service/internal/consul"
	"product-service/internal/products"
	"product-service/pkg/ctxmanage"
	"product-service/pkg/logkey"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) addToCart(c *gin.Context) {

	traceId := ctxmanage.GetTraceIdOfRequest(c)

	if c.Request.ContentLength > 5*1024 {
		// Log an error indicating that the request body size limit was breached.
		slog.Error("request body limit breached",
			slog.String(logkey.TraceID, traceId),
			slog.Int64("Size Received", c.Request.ContentLength),
		)

		// Respond with HTTP 400 Bad Request and an appropriate error message.
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "payload exceeding size limit",
		})
		return
	}

	claims, err := ctxmanage.GetAuthClaimsFromContext(c.Request.Context())
	if err != nil {
		slog.Error(
			"missing claims",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User Id is required"})
		return
	}

	userId := claims.Subject

	var newCart products.NewCartLine

	err = c.ShouldBindBodyWithJSON(&newCart)

	if err != nil {
		// Log an error if JSON parsing or validation fails, along with the trace ID.
		slog.Error("json validation error",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 400 Bad Request and indicate the error.
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": http.StatusText(http.StatusBadRequest),
		})
		return
	}

	err = h.validate.Struct(newCart)
	if err != nil {
		// Log an error if validation fails, along with the trace ID and specific error message.
		slog.Error("validation failed",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 400 Bad Request and indicate the need for correct input formats.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide values in correct format",
		})
		return
	}

	ctx := c.Request.Context()

	// Attempt to insert the new user into the database using the `InsertUser` method.
	err = h.p.InsertOrUpdateCart(ctx, userId, newCart)
	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in creating the product",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Add To Cart Failed",
		})
		return
	}
	// Respond with HTTP 200 OK and return the created user's data as JSON.
	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart"})

}

func (h *Handler) checkout(c *gin.Context) {

	traceId := ctxmanage.GetTraceIdOfRequest(c)

	if c.Request.ContentLength > 5*1024 {
		// Log an error indicating that the request body size limit was breached.
		slog.Error("request body limit breached",
			slog.String(logkey.TraceID, traceId),
			slog.Int64("Size Received", c.Request.ContentLength),
		)

		// Respond with HTTP 400 Bad Request and an appropriate error message.
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "payload exceeding size limit",
		})
		return
	}

	claims, err := ctxmanage.GetAuthClaimsFromContext(c.Request.Context())
	if err != nil {
		slog.Error(
			"missing claims",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User Id is required"})
		return
	}

	// FETCH USER ID

	userId := claims.Subject

	ctx := c.Request.Context()

	// Check if items are inProgress for checkout
	// Attempt to insert the new user into the database using the `InsertUser` method.
	cartRet, err := h.p.FetchCartItems(ctx, userId, products.StatusInProgress)
	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in fetching the cart",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "fetch Cart Failed",
		})
		return
	}

	if len(cartRet.LineItems) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "There are no items in the cart ",
		})
		return
	}

	var orderReq products.OrderRequest
	orderReq.LineItems = cartRet.LineItems

	orderId := cartRet.OrderId

	//UPDATE THE STATUS TO PENDING

	err = h.p.UpdateCartStatusFromInProgressToPending(ctx, userId)

	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in updating the cart status",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "fetch Cart Failed",
		})
		return
	}

	type OrderServiceResponse struct {
		CheckoutSessionID string `json:"checkout_session_id"`
	}

	//caLL ORDER SERVICE CHECKOUT
	orderChan := make(chan OrderServiceResponse, 1) // For customer ID

	go func() {

		requestBody, err := json.Marshal(orderReq)
		if err != nil {
			slog.Error("error marshalling request body",
				slog.String(logkey.TraceID, traceId),
				slog.Any("error", err))
			orderChan <- OrderServiceResponse{}
			return
		}
		address, port, err := consul.GetServiceAddress(h.client, "orders")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			orderChan <- OrderServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/orders/cartcheckout/v2/%s", address, port, orderId)
		slog.Info("httpQuery: "+httpQuery, slog.String(logkey.TraceID, traceId))

		// Make the HTTP POST request
		fmt.Println(string(requestBody))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, httpQuery, bytes.NewBuffer(requestBody))
		if err != nil {
			slog.Error("error creating request", slog.String(logkey.TraceID, traceId), slog.Any("error", err.Error()))
			orderChan <- OrderServiceResponse{}
			return
		}
		authorizationHeader := c.Request.Header.Get("Authorization")
		req.Header.Set("Authorization", authorizationHeader)

		resp, err := http.DefaultClient.Do(req)
		fmt.Println(resp)
		if err != nil {
			slog.Error("error fetching order service", slog.String(logkey.TraceID, traceId))
			orderChan <- OrderServiceResponse{}
			return
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching checkout session from order service", slog.String(logkey.TraceID, traceId))
			orderChan <- OrderServiceResponse{}
			return
		}

		defer resp.Body.Close()

		var orderServiceResponse OrderServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&orderServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			orderChan <- OrderServiceResponse{}
			return
		}
		// Print the customer Id if fetched successfully
		slog.Info("successfully fetched stripe customer id", slog.String(logkey.TraceID, traceId))
		orderChan <- orderServiceResponse
	}()

	orderServiceResponse := <-orderChan
	if orderServiceResponse.CheckoutSessionID == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching checkout session"})
		return
	}

	c.JSON(http.StatusOK, orderServiceResponse)

}

func (h *Handler) fetchCartDetails(c *gin.Context) {

	traceId := ctxmanage.GetTraceIdOfRequest(c)

	if c.Request.ContentLength > 5*1024 {
		// Log an error indicating that the request body size limit was breached.
		slog.Error("request body limit breached",
			slog.String(logkey.TraceID, traceId),
			slog.Int64("Size Received", c.Request.ContentLength),
		)

		// Respond with HTTP 400 Bad Request and an appropriate error message.
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "payload exceeding size limit",
		})
		return
	}

	claims, err := ctxmanage.GetAuthClaimsFromContext(c.Request.Context())
	if err != nil {
		slog.Error(
			"missing claims",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User Id is required"})
		return
	}

	// FETCH USER ID

	userId := claims.Subject

	ctx := c.Request.Context()

	// Check if items are inProgress for checkout
	// Attempt to insert the new user into the database using the `InsertUser` method.
	cartRet, err := h.p.FetchCartDetails(ctx, userId, products.StatusInProgress)
	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in fetching the cart",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "fetch Cart Failed",
		})
		return
	}
	if len(cartRet) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "Empty Cart",
		})
		return
	}

	c.JSON(http.StatusOK, cartRet)

}

func (h *Handler) deleteCartByID(c *gin.Context) {
	traceId := ctxmanage.GetTraceIdOfRequest(c)

	// Ensure the request body isn't too large
	if c.Request.ContentLength > 5*1024 {
		slog.Error("request body limit breached",
			slog.String(logkey.TraceID, traceId),
			slog.Int64("Size Received", c.Request.ContentLength),
		)

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "payload exceeding size limit",
		})
		return
	}

	// Extract claims and validate user authentication
	claims, err := ctxmanage.GetAuthClaimsFromContext(c.Request.Context())
	if err != nil {
		slog.Error("missing claims",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "User Id is required"})
		return
	}

	// Extract user ID from claims
	userId := claims.Subject

	// Extract cart ID from request URL parameter
	cartID := c.Param("id")
	if cartID == "" {
		slog.Error("missing cart ID",
			slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Cart ID is required",
		})
		return
	}

	ctx := c.Request.Context()

	// Call the service layer to delete the cart item
	err = h.p.DeleteCartByIDIfPending(ctx, cartID, userId)
	if err != nil {
		// Handle specific errors (e.g., no rows affected)
		if errors.Is(err, sql.ErrNoRows) {
			slog.Warn("no cart item found to delete",
				slog.String(logkey.TraceID, traceId),
				slog.String("cartID", cartID),
			)

			c.JSON(http.StatusNotFound, gin.H{
				"message": "Cart item not found",
			})
			return
		}

		// Log any unexpected errors
		slog.Error("failed to delete cart item",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete cart item",
		})
		return
	}

	// Respond with success
	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item successfully deleted",
	})
}
