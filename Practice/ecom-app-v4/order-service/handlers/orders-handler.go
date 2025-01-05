package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"order-service/consul"
	"order-service/internal/auth"
	"order-service/pkg/ctxmanage"
	"order-service/pkg/logkey"
	"order-service/protohandler"
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
	//c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
	userId := claims.Subject
	ctx := c.Request.Context()
	err = h.o.CreateOrder(ctx, orderId, userId, productID, sessionStripe.AmountTotal)
	if err != nil {
		slog.Error("error creating order", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order"})
		return
	}
	//hit server here
	protoresp, err := protohandler.HitServer(h.protoclient, productID)
	if err != nil {
		slog.Error("error with grpc server", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hit grpc"})
		return
	}
	// Respond with the Stripe session ID
	c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL, "protoresp": protoresp})
}

func (h *Handler) CheckoutWithGrpc(c *gin.Context) {
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
		//hit server here
		protoresp, err := protohandler.HitServer(h.protoclient, productID)
		if err != nil {
			slog.Error("error with grpc server", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hit grpc"})
			productChan <- ProductServiceResponse{}
			return
		}

		slog.Info("successfully hit grpc and returned", slog.String(logkey.TraceID, traceId))
		pr := protoresp.GetProdOrder()

		slog.Info("successfully hit grpc and returned", slog.String(logkey.TraceID, traceId), slog.String("product data", fmt.Sprintf("pr", pr)))

		fmt.Println(int(pr.GetStock()), pr.GetPriceId())
		productServiceResponse := ProductServiceResponse{ProductID: productID, Stock: int(pr.GetStock()), PriceID: pr.GetPriceId()}

		productChan <- productServiceResponse
	}()

	// Set a timeout (e.g., 10 seconds)
	//TODO remove and test out why product channel is not benig picked up
	timeout := time.After(100 * time.Second)

	select {

	case <-timeout:
		// Handle the timeout case
		slog.Error("Timeout waiting for product service response", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusRequestTimeout, gin.H{"message": "Request to product service timed out"})
		return
	}

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
	//c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
	userId := claims.Subject
	ctx := c.Request.Context()
	err = h.o.CreateOrder(ctx, orderId, userId, productID, sessionStripe.AmountTotal)
	if err != nil {
		slog.Error("error creating order", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order"})
		return
	}
	// Respond with the Stripe session ID
	c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
}

func (h *Handler) CartCheckout(c *gin.Context) {
	//TODO: Add the order in the orders table, and mark that as pending

	// Get the traceId from the request for tracking logs
	traceId := ctxmanage.GetTraceIdOfRequest(c)
	claims, ok := c.Request.Context().Value(auth.ClaimsKey).(auth.Claims)
	if !ok {
		slog.Error("claims not found", slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusUnauthorized})
		return
	}

	var req CartOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error(
			"invalid request body",
			slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Invalid request body"})
		return
	}

	orderId := c.Param("orderId")
	if orderId == "" {
		orderId = uuid.NewString()
	}

	// Validate that the list is not empty
	if len(req.LineItems) < 1 {
		slog.Error(
			"empty lineitems",
			slog.String(logkey.TraceID, traceId))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Line Items cannot be empty"})
		return
	}

	//conver to a map
	// Convert to map
	productMap := make(map[string]LineItem)

	var productIds []string

	// Populate the map
	for _, item := range req.LineItems {
		productMap[item.ProductId] = item
		productIds = append(productIds, item.ProductId)
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

	productChan := make(chan []ProductServiceResponse, 1) // For stock and price information
	go func() {
		address, port, err := consul.GetServiceAddress(h.client, "products")
		if err != nil {
			slog.Error("service unavailable", slog.String(logkey.TraceID, traceId),
				slog.String(logkey.ERROR, err.Error()))
			productChan <- nil
			return
		}
		// Create a JSON body with the list of product IDs
		fmt.Println("productsId", productIds)
		fmt.Println(map[string][]string{"productIds": productIds})
		requestBody, err := json.Marshal(map[string][]string{"productIds": productIds})
		if err != nil {
			slog.Error("error marshalling request body",
				slog.String(logkey.TraceID, traceId),
				slog.Any("error", err))
			productChan <- nil
			return
		}

		httpURL := fmt.Sprintf("http://%s:%d/products/stock", address, port)

		// Make the HTTP POST request
		fmt.Println(string(requestBody))
		resp, err := http.Post(httpURL, "application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			slog.Error("error fetching product service",
				slog.String(logkey.TraceID, traceId),
				slog.Any("error", err))
			productChan <- nil
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			slog.Error("error fetching product information", slog.String(logkey.TraceID, traceId))
			productChan <- nil
			return
		}
		var productServiceResponse []ProductServiceResponse
		err = json.NewDecoder(resp.Body).Decode(&productServiceResponse)
		if err != nil {
			slog.Error("error binding json", slog.String(logkey.TraceID, traceId), slog.Any(logkey.ERROR, err.Error()))
			productChan <- nil
			return
		}
		fmt.Println("product reso", productServiceResponse)

		productChan <- productServiceResponse
	}()

	userServiceResponse := <-userChan
	if userServiceResponse.StripCustomerId == "" {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching stripe customer id"})
		return
	}
	stockData := <-productChan

	if stockData == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching product information"})
		return
	}
	var lineItems []*stripe.CheckoutSessionLineItemParams

	for _, stockVal := range stockData {
		priceID := stockVal.PriceID
		stock := stockVal.Stock
		productID := stockVal.ProductID
		if stock <= 0 || priceID == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error fetching product information"})
			return
		}
		//Create line items
		lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
			Price: stripe.String(stockVal.PriceID),
			//Todo make this dynamic
			//Quantity: stripe.Int64(stockVal.Quantity),
			Quantity: stripe.Int64(int64(productMap[productID].Quantity)),
		})
		//create metadata
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

	// Convert struct to JSON string
	reqJSON, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("Error marshaling struct: %v\n", err)
		return
	}
	// Proceed to create Stripe checkout session
	// var lineItems []*stripe.CheckoutSessionLineItemParams
	// for _, product := range products {
	// 	lineItems = append(lineItems, &stripe.CheckoutSessionLineItemParams{
	// 		Price:    stripe.String(product.PriceID),
	// 		Quantity: stripe.Int64(product.Quantity),
	// 	})
	// }
	//TODO SEND LIST OF PRODUCT IDS
	params := &stripe.CheckoutSessionParams{
		Customer:                 stripe.String(userServiceResponse.StripCustomerId),
		SubmitType:               stripe.String("pay"),
		Currency:                 stripe.String(string(stripe.CurrencyINR)),
		BillingAddressCollection: stripe.String("auto"),
		LineItems:                lineItems,
		Mode:                     stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL:               stripe.String("https://example.com/success"),
		//ExpiresAt:
		CancelURL: stripe.String("https://example.com/cancel"),
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			Metadata: map[string]string{
				"order_id": orderId,
				"user_id":  claims.Subject, // userID in jwt token
				"products": string(reqJSON),
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
	//TODO change string(reqJSON) to prod id
	slog.Info("successfully initiated Stripe checkout session", slog.String("Trace ID", traceId), slog.String("ProductID", string(reqJSON)), slog.String("CheckoutSessionID", sessionStripe.ID))

	slog.Info("successfully initiated Stripe checkout session", slog.String("CheckoutSessionID", sessionStripe.URL))

	// Respond with the Stripe session ID
	//c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
	userId := claims.Subject
	ctx := c.Request.Context()
	//TODO chnge string(reqJSON) for productId
	fmt.Println(orderId, userId, productIds[0], sessionStripe.AmountTotal)
	err = h.o.CreateOrder(ctx, orderId, userId, productIds[0], sessionStripe.AmountTotal)
	if err != nil {
		fmt.Println(err)
		slog.Error("error creating order", slog.String(logkey.TraceID, traceId), slog.String(logkey.ERROR, err.Error()))
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create order"})
		return
	}

	// Respond with the Stripe session ID
	c.JSON(http.StatusOK, gin.H{"checkout_session_id": sessionStripe.URL})
}

/*
*
                  +---------------------+
                  |       START         |
                  |   API Call Starts   |
                  +---------------------+
                            |
                            v
                +----------------------+
                |   Check STRIPE Key   |
                | (Environment Config) |
                +----------------------+
                            |
           STRIPE Key FOUND | STRIPE Key MISSING
                  |                   v
                  v       +----------------------------+
    +------------------+  | Respond with Error:         |
    | Extract User     |  | "Stripe test key not found" |
    | Claims & TraceID |  +----------------------------+
    +------------------+
                  |
                  v
          +-------------------+
          | Extract ProductID |
          |   From Request    |
          +-------------------+
                  |
        ProductID FOUND | ProductID MISSING
                  |                   v
                  v       +--------------------------------+
    +----------------------------------+| Respond with Error: |
    | Create Channels for Concurrent  || "Product ID Missing"|
    |     Service Calls               |+--------------------------------+
    +----------------------------------+
                  |
                  v
   +---------------------------------------+
   | Start Parallel Service Calls          |
   | 1. Call User Service (Stripe ID)      |
   | 2. Call Product Service (Stock/Price) |
   +---------------------------------------+
                  |
          +-------------------+  +-------------------+
          | Wait for User ID  |  | Wait for Product   |
          +-------------------+  | Details            |
                  |               +-------------------+
       Stripe ID FOUND |   Product Details FOUND
                  |                    |
                  v                    v
       +---------------------------------------+
       | Validate Results:                    |
       | - Valid Stripe Customer ID           |
       | - Valid Product Details (Stock > 0,  |
       |   PriceID Exists)                    |
       +---------------------------------------+
                  |
          Validation PASSED | Validation FAILED
                  |                   	v
                  v      				----------------------------+
       +----------------------------+ 	Respond with Error    |
       | Create Stripe Checkout     | 	"Invalid Inputs"      |
       | Session with User & Product|	-----------------------+
       +----------------------------+
                  |// Create the order in the orders table with a pending status
                  v
     +--------------------------------+
     |  Respond with Checkout URL     |
     |  (Stripe Session Created)      |
     +--------------------------------+
                  |
                  v
           +-------------+
           |     END     |
           +-------------+
*/
