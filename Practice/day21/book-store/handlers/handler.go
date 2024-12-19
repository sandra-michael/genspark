package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func SetupGINRoutes() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", Ping)
	return r
}

func HandleBookStore() *gin.Engine {

	// Step 1: Initialize OpenTelemetry
	traceProvider, err := initOpenTelemetry()
	if err != nil {
		panic(err)
	}
	defer traceProvider.Shutdown(context.Background())

	r := gin.Default()

	// Step 3: Add OpenTelemetry middleware to Gin
	// This will automatically trace all incoming HTTP requests handled by Gin.
	//needs to be the first middleware which is being used
	r.Use(otelgin.Middleware("book-micro"))

	v1 := r.Group("v1/books")
	{
		v1.POST("/createTable", CreateTable)
		v1.POST("/", Insertfunc)
	}
	return r
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
			semconv.ServiceNameKey.String("book-micro"), // Set the service name for tracing
		)),
	)

	// Register the global TracerProvider and propagators
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})
	return traceProvider, nil
}