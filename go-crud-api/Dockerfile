# Step 1: Build stage
FROM golang:1.23-alpine AS build

# Set environment variables
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

# Set working directory inside the container
WORKDIR /app

# Copy the Go modules files to download dependencies first
COPY go.mod go.sum ./

# Download Go modules dependencies
RUN go mod download

# Copy the source code to the container
COPY . .

# Build the Go Gin app
RUN go build -o /go-crud-api

# Step 2: Runtime stage
FROM alpine:latest

# Create a non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy the compiled binary from the build stage
COPY --from=build /go-crud-api /app/go-crud-api

# Set user to non-root for security
USER appuser

# Expose the port that the Gin app uses (usually 8080)
EXPOSE 8080

# Command to run the app
CMD ["/app/go-crud-api"]
