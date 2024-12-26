# Start with a base image containing Go
FROM golang:1.22 AS builder

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN go mod download

# Build the Go application
RUN GOOS=linux GOARCH=arm64 go build -o main ./cmd/waiter

# Use a minimal base image for production
FROM alpine:latest

# Set up working directory and copy binary
WORKDIR /root/
COPY --from=builder /app/main .

RUN apk add file
RUN file main

# Expose application port
EXPOSE 8080

# Run the application
RUN chmod +x main
CMD ["./main"]
