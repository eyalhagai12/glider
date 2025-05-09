# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download -x

# Copy the source code and build the application
COPY . .
RUN go build -o main main.go

# Stage 2: Create a lightweight image to run the application
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the binary
CMD ["./main"]

