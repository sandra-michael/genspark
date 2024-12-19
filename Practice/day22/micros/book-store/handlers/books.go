package handlers

import (
	"book-store/models"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func Ping(c *gin.Context) {
	time.Sleep(10 * time.Second)
	c.JSON(200, gin.H{"message": "pong"})
}

func CreateTable(c *gin.Context) {

	conn, err := models.NewConn()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	err = conn.CreateBookTable(c.Request.Context())
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusCreated, "Book Table")

}

func Insertfunc(c *gin.Context) {
	tracer := otel.Tracer("book-micro")
	ctx, span := tracer.Start(c.Request.Context(), "Insertfunc")
	defer span.End()

	var newBook models.NewBook

	traceId := span.SpanContext().TraceID().String()

	// Call BindJSON to bind the received JSON to
	// newBooks.
	if err := c.BindJSON(&newBook); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Converting json to struct error "})
		return
	}
	conn, err := models.NewConn()
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error forming a new connection"})
		return
	}

	book, err := conn.InsertBook(ctx, newBook)
	if err != nil {
		fmt.Println(err)
		// Handle and record any errors in the span
		span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusBadRequest)) // HTTP 400
		//span.SetAttributes(attribute.String("user_id", userId))                          // Attach user ID
		span.SetAttributes(attribute.String("traceId", traceId))
		span.AddEvent("UNABLE TO INSERT BOOK") // Record event in tracing span// Attach trace ID
		span.SetStatus(codes.Error, err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error during insert "})
		return
	}
	span.SetAttributes(semconv.HTTPResponseStatusCodeKey.Int(http.StatusOK)) // HTTP 200

	client := &http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}

	// The context `ctx` contains tracing metadata that will be sent with the request.
	url := "http://localhost:8083/logbook/" + strconv.Itoa(book.ID)
	log.Println("url", url)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		// If constructing the HTTP request fails (e.g., incorrect URL or method), log the error
		// and respond to the client with a 500 Internal Server Error.
		log.Printf("Failed to construct request for the order service: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 4: Inject trace context metadata into the HTTP request headers.
	// OpenTelemetry uses this to propagate trace data to the downstream Order Service.
	// This ensures distributed tracing works seamlessly across multiple services.
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Step 5: Execute the HTTP request to the Order Service using the instrumented client.
	resp, err := client.Do(req)
	if err != nil {
		// If there's an error communicating with the Order Service (e.g., server down or network issue),
		// log the error and respond to the client with a 500 Internal Server Error.
		log.Printf("Failed to call the order service: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 6: Read the response body returned by the Order Service.
	// The `io.ReadAll` function reads the entire body into memory as a byte slice.
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		// If reading the response body fails, log the error and respond with a 500 error to the client.
		log.Printf("Failed to read response from order service: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Step 7: Set the status of the tracing span to "Ok" to indicate that
	// the request to the Order Service was successful and the response was processed correctly.
	span.SetStatus(codes.Ok, "order service response received")

	// Step 8: Log the response received from the external service for debugging or monitoring purposes.
	// This provides visibility into what the downstream service returned.
	log.Printf("Order service response: %s", string(b))

	// Step 9: Send the response from the Order Service back to the client.
	// Respond with HTTP status 200 (OK) if everything is successful.
	// Use `c.String` to return the response as a plain-text string.
	c.String(http.StatusOK, string(b))

	//c.IndentedJSON(http.StatusCreated, book)
	return

}
