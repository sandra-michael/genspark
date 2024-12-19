# micro-1/Dockerfile

# Stage 1: Build the Go application
FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app



# Set the target OS and architecture // very important if deploying in linux containers in cloud
# ENV GOOS=linux
# ENV GOARCH=amd64

# Copy go.mod and go.sum first to leverage layer caching
COPY go.mod go.sum ./
#if copy go.mod and go.sum before then run go mod tidy and download immediately
# it would cache the layer so both command would only run if there any changes in go.mod or go.sum
# Download and cache dependencies
RUN go mod download

# Run go mod tidy to ensure that dependencies are consistent
RUN go mod tidy

# Copy the source code into the container
COPY . .

## Download and cache dependencies
#RUN go mod download
## Run go mod tidy to ensure that dependencies are consistent
#RUN go mod tidy

# Build the Go application
#RUN go build -o hello-service .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway-service .

# Stage 2: Create a lightweight runtime image

FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the previous stage
COPY --from=builder /app/gateway-service .


EXPOSE 80
# Command to run the application
CMD ["./gateway-service"]