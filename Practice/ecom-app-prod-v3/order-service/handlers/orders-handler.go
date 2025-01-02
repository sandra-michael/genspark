package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"order-service/consul"
	"order-service/internal/auth"
	"order-service/pkg/ctxmanage"
	"order-service/pkg/logkey"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/checkout/session"
)

func (h *Handler) Checkout(c *gin.Context) {
	//TODO: Add the order in the orders table, and mark that as pending

	// Get the traceId from the request for tracking logs
	traceId := ctxmanage.GetTraceIdOfRequest(c)
	claims, ok := c.Request.Context().Value(auth.ClaimsKey).(auth.Claims)
	if !ok {
		slog.Error("claims not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusUnauthorized})
		return
	}

	type UserServiceResponse struct {
		StripCustomerId string `json:"stripe_customer_id"`
	}
	type ProductServiceResponse struct {
		ProductID string `json:"product_id"`
		Stock     int    `json:"stock"`
		PriceID   string `json:"price_id"`
	}

	productID := c.Param("productID")
	if productID == "" {
		slog.Error("missing product id", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
		return
	}

	// Create channels for goroutine results
	userChan := make(chan UserServiceResponse, 1) // For customer ID

	if h.client == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "consul client is not initialized"})
	}

	go func() {

		address, port, err := consul.GetServiceAddress(h.client, "users")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/users/stripe", address, port)
		slog.Info("httpQuery: "+httpQuery, slog.String(logkey.TraceID, traceId))
		ctx, cancel := context.WithTimeout(c.Request.Context(), 50*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpQuery, nil)
		if err != nil {
			slog.Error("error creating request", slog.String(logkey.TraceID, traceId), slog.Any("error", err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		authorizationHeader := c.Request.Header.Get("Authorization")
		req.Header.Set("Authorization", authorizationHeader)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			slog.Error("error fetching user service", slog.String(logkey.TraceID, traceId))
			userChan <- UserServiceResponse{}
			return
		}
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching stripe id from user service", slog.String(logkey.TraceID, traceId))
			userChan <- UserServiceResponse{}
			return
		}

		defer resp.Body.Close()

		var userServiceResponse UserServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&userServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			userChan <- UserServiceResponse{}
			return
		}
		// Print the customer Id if fetched successfully
		slog.Info("successfully fetched stripe customer id", slog.String(logkey.TraceID, traceId))
		userChan <- userServiceResponse
	}()

	productChan := make(chan ProductServiceResponse, 1) // For stock and price information
	go func() {
		address, port, err := consul.GetServiceAddress(h.client, "products")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			productChan <- ProductServiceResponse{}
			return
		}
		httpQuery := fmt.Sprintf("http://%s:%d/products/stock/%s", address, port, productID)
		resp, err := http.Get(httpQuery)
		if err != nil {
			slog.Error("error fetching product service", slog.String(logkey.TraceID, traceId))
			productChan <- ProductServiceResponse{}
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching product information", slog.String(logkey.TraceID, traceId))
			productChan <- ProductServiceResponse{}
			return
		}
		var productServiceResponse ProductServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&productServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			productChan <- ProductServiceResponse{}
			return
		}
		productChan <- productServiceResponse
	}()

	userServiceResponse := <-userChan
	if userServiceResponse.StripCustomerId == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching stripe customer id"})
		return
	}
	stockPriceData := <-productChan
	priceID := stockPriceData.PriceID
	stock := stockPriceData.Stock
	if stock <= 0 || priceID == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching product information"})
		return
	}
	//c.JSON(http.StatusOK, gin.H{"customerId": userServiceResponse.StripCustomerId, "price_id": priceID, "stock": stock})

	// Step 1: Retrieve the Stripe secret key from the environment variables
	sKey := os.Getenv("STRIPE_TEST_KEY")
	if sKey == "" {
		slog.Error("Stripe secret key not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Stripe secret key not found"})
	}

	// Step 2: Assign the Stripe API key to the Stripe library's internal configuration
	stripe.Key = sKey
	orderId := uuid.NewString()
	// Proceed to create Stripe checkout session
	params := &stripe.CheckoutSessionParams{
		Customer:                 stripe.String(userServiceResponse.StripCustomerId),
		SubmitType:               stripe.String("pay"),
		Currency:                 stripe.String(string(stripe.CurrencyINR)),
		BillingAddressCollection: stripe.String("auto"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1), // Adjust quantity as needed
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String("https://example.com/success"),
		//ExpiresAt:
		CancelURL: stripe.String("https://example.com/cancel"),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"order_id":   orderId,
				"user_id":    claims.Subject, // userID in jwt token
				"product_id": productID,
			},
		},
	}

	sessionStripe, err := session.New(params)
	if err != nil {
		slog.Error("error creating Stripe checkout session", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create Stripe checkout session"})
		return
	}

	// Log success operation
	slog.Info("successfully initiated Stripe checkout session", slog.String("Trace ID", traceId), slog.String("ProductID", productID), slog.String("CheckoutSessionID", sessionStripe.ID))

	// Respond with the Stripe session ID
	c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
}
