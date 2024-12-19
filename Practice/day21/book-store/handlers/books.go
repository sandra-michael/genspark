package handlers

import (
	"book-store/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
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

	c.IndentedJSON(http.StatusCreated, book)
	return

}
