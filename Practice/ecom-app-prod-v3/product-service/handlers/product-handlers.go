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
		slog.Error("error in creating the user",
			slog.String(logkey.TraceID, traceId),
			slog.String(logkey.ERROR, err.Error()),
		)

		// Respond with HTTP 500 Internal Server Error indicating that user creation failed.
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "User Creation Failed",
		})
		return
	}

	h.p.CreatePricingStripe(ctx, product.ID, paise, product.Name)
	// Respond with HTTP 200 OK and return the created user's data as JSON.
	c.JSON(http.StatusOK, product)

}

// func validatePrice(price string) (uint, error) {
// 	//trim
// 	price = strings.Trim(price, " ")

// 	pr := strings.Split(price, ".")

// 	//check the length of split
// 	if len(pr) > 2 {
// 		return 0, fmt.Errorf("THe provided price is invalid")
// 	}

// 	//converting part 1 to uint to avoid negative
// 	rupee, err := strconv.ParseUint(pr[0], 10, 64)

// 	if err != nil {
// 		return 0, fmt.Errorf("THe provided price needs to be valid before decimal")
// 	}

// 	paisa := 0
// 	if len(price) == 2 {
// 		paisa, err = strconv.ParseUint(pr[1], 10, 64)
// 		if err != nil || paisa > 99 {
// 			return 0, fmt.Errorf("invalid price, please provide price in valid format")
// 		}

// 		// append 0 if paisa part has only one digit
// 		// e.g INR 99.2 => Convert it to 9900 + 20 = 9920
// 		if paisa < 10 {
// 			paisa *= 10
// 		}
// 	}
// 	return rupee*100 + paisa, nil

// }

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
