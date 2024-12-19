package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

// q. From book store call another microservice that print a book has been added in the database
//     and return a response from this new microservice: Book:id logging done
//     Add distributed tracing to this

func main() {
	// Step 1: Initialize OpenTelemetry
	traceProvider, err := initOpenTelemetry()
	if err != nil {
		panic(err)
	}
	defer traceProvider.Shutdown(context.Background())

	// Step 2: Create a Gin router
	r := gin.Default()

	// Step 3: Add OpenTelemetry middleware to Gin
	r.Use(otelgin.Middleware("print-micro"))

	// Step 4: Define the  endpoint
	r.GET("/logbook/:bookid", printMicro)

	// Step 5: Start the server on port 8083
	if err := r.Run(":8083"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func printMicro(c *gin.Context) {

	//Extract trace from incoming request usin propagator
	propagator := otel.GetTextMapPropagator()

	//Extracting context from the propagatot
	extractedCtx := propagator.Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

	log.Printf("viewing headers")
	for k, v := range c.Request.Header {
		log.Println(k, v) // Log each header key and its value(s).
	}

	//creating a ctx or span
	//preferred  prctice span name = func name
	_, span := otel.Tracer("print-micro").Start(extractedCtx, "Print Handler")
	defer span.End()

	bookId := c.Param("bookid")
	message := "Book: " + bookId + " logging done "
	c.String(http.StatusOK, message)
	log.Println(message)

}

func initOpenTelemetry() (*trace.TracerProvider, error) {
	// Set up the OTLP trace exporter to send tracing data to the OpenTelemetry Collector
	traceExporter, err := otlptracehttp.New(
		context.Background(),
		otlptracehttp.WithInsecure(),                 // No TLS for local development
		otlptracehttp.WithEndpoint("localhost:4318"), // Collector/Jaeger endpoint
	)
	if err != nil {
		return nil, err
	}

	// Configure a TracerProvider
	traceProvider := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()), // Sample all traces
		trace.WithBatcher(traceExporter),        // Batch traces in export
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("print-micro"), // Set the service name for tracing
		)),
	)

	// Register the global TracerProvider and propagators
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return traceProvider, nil
}
