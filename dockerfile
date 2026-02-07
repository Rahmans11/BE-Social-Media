# Stage 1: Build the Go application
FROM golang:alpine AS builder

# Install git and other build dependencies
RUN apk update && apk add --no-cache git build-base

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to leverage Docker's cache
COPY go.mod go.sum ./

# Download all the dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application binary. Adjust the path to your main package if needed (e.g., ./cmd/server/)
RUN go build -o main ./cmd/main.go

# Stage 2: Run the application
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your Gin API listens on (default is usually 8080)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]