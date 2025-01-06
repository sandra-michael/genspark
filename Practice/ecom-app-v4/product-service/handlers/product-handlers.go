package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"product-service/internal/products"
	"product-service/pkg/ctxmanage"
	"product-service/pkg/logkey"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createProduct(c *gin.Context) {

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

	var newProduct products.NewProduct

	err := c.ShouldBindBodyWithJSON(&newProduct)

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

	err = h.validate.Struct(newProduct)
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

	//TODO ADD VALIDATION FOR PRICE
	paise, err := RupeesToPaise(newProduct.Price)

	if err != nil {
		// Log an error if validation fails, along with the trace ID and specific error message.
		slog.Error("validation failed",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 400 Bad Request and indicate the need for correct input formats.
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "please provide Price in correct format",
		})
		return
	}
	newProduct.Price = strconv.Itoa(int(paise))

	// Extract the context from the HTTP request to pass it to the service layer.
	ctx := c.Request.Context()

	// Attempt to insert the new user into the database using the `InsertUser` method.
	product, err := h.p.InsertProduct(ctx, newProduct)
	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in creating the product",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Product Creation Failed",
		})
		return
	}

	h.p.CreatePricingStripe(ctx, product.ID, paise, product.Name)
	// Respond with HTTP 200 OK and return the created user's data as JSON.
	c.JSON(http.StatusOK, product)

}

func RupeesToPaise(priceStr string) (uint64, error) {
	fmt.Println("Input price:", priceStr)
	//trim extra space from price
	priceStr = strings.Trim(priceStr, " ")

	//split the price based by dot(.)
	prices := strings.Split(priceStr, ".")
	var rupee, paisa uint64
	if len(prices) == 0 || len(prices) > 2 {
		return 0, fmt.Errorf("invalid price, empty price field or more than one dot(.)")
	}

	rupee, err := strconv.ParseUint(prices[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price, not a valid number")
	}

	if len(prices) == 2 {

		if len(prices[1]) > 2 {
			return 0, fmt.Errorf("invalid price, please provide price in valid format")
		}
		paisa, err = strconv.ParseUint(prices[1], 10, 64)
		if err != nil || paisa > 99 {
			return 0, fmt.Errorf("invalid price, please provide price in valid format")
		}

		// append 0 if paisa part has only one digit
		// e.g INR 99.2 => Convert it to 9900 + 20 = 9920
		if paisa < 10 {
			paisa *= 10
		}
	}
	return rupee*100 + paisa, nil
}

// fetch all products
func (h *Handler) fetchAllProducts(c *gin.Context) {

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

	// Extract the context from the HTTP request to pass it to the service layer.
	ctx := c.Request.Context()
	// Attempt to insert the new user into the database using the `InsertUser` method.
	products, err := h.p.FetchAllProducts(ctx)
	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in fetching the products",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Product Fetch Failed",
		})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"message": "There are no products",
		})
		return

	}

	// Respond with HTTP 200 OK and return the created user's data as JSON.
	c.JSON(http.StatusOK, products)
}

/*
*
We cannot update the price value in stripe
https://stackoverflow.com/questions/76179501/getting-this-error-received-unknown-parameter-unit-amount-when-updating-produc

You can't update the unit_amount of a price. You can achieve it (by setting its active to false) and create a new price with the new unit_amount.

# Refer the to API reference for more details

# Will add this step in TODO since it is complicated and not allow price updates

*
*/
func (h *Handler) updateProduct(c *gin.Context) {

	traceId := ctxmanage.GetTraceIdOfRequest(c)

	// Extract the product ID from the URL parameters
	productID := c.Param("productID")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product ID is required"})
		return
	}

	// Parse the request body
	var req products.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request body", "error": err.Error()})
		return
	}

	err := h.validate.Struct(req)
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
	/**
	Commenting out this peice of code due to the complexity od price update with stripe
	will add to TODO LIST
	**/
	// var paise uint64
	// fmt.Println("Req price value ", req.Price, req.Price != "")
	// if req.Price != "" {
	// 	//TODO ADD VALIDATION FOR PRICE
	// 	paiseRet, err := RupeesToPaise(req.Price)

	// 	if err != nil {
	// 		// Log an error if validation fails, along with the trace ID and specific error message.
	// 		slog.Error("validation failed",
	// 			slog.String(logkey.TraceID, traceId),
	// 			slog.String(logkey.ERROR, err.Error()),
	// 		)

	// 		// Respond with HTTP 400 Bad Request and indicate the need for correct input formats.
	// 		c.JSON(http.StatusBadRequest, gin.H{
	// 			"error": "please provide Price in correct format",
	// 		})
	// 		return
	// 	}
	// 	req.Price = strconv.Itoa(int(paiseRet))
	// 	/**
	// 	Problem:
	// 	Inside the if block, paise := 123 creates a new local variable paise, and it doesn’t modify the outer paise variable.
	// 	The outer paise remains at its initial value (0 since it's declared as uint64 and default initialized to 0).
	// 	Solution:
	// 	You need to avoid shadowing the paise variable by directly modifying the outer paise variable inside the if block.
	// 	**/
	// 	paise = paiseRet
	// 	fmt.Println("paise within func ", paise)

	// }
	ctx := c.Request.Context()

	// Attempt to insert the new user into the database using the `InsertUser` method.
	err = h.p.UpdateProduct(ctx, productID, req)
	if err != nil {
		// Log an error if user creation fails, along with the trace ID and specific error message.
		slog.Error("error in creating the product",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Ptoduct Update Failed",
		})
		return
	}

	/**
		TODO when price update with stripe is sorted
	**/
	// fmt.Println("before sending to stripe", paise)

	// //TODO CAN ALSO UPDATE FOR PRODUCT NAME
	// //if price is being updated hit stripe
	// if req.Price != "" {
	// 	h.p.UpdatePricingStripe(ctx, productID, paise)
	// }
	// Respond with success
	c.JSON(http.StatusOK, gin.H{"message": "Updated Product Details"})
}
