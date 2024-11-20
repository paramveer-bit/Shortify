# Use Golang image for the backend
FROM golang:1.23-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Install dependencies (this will use the cache if go.mod and go.sum are unchanged)
RUN go mod tidy

# Copy the rest of the Go files
COPY . .

# Build the Go app
RUN go build -o main .

# Create a smaller image to run the Go app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the Go executable from the builder
COPY --from=builder /app/main .

# Expose the backend port (e.g., 8080)
EXPOSE 8080

# Run the Go app
CMD ["./main"]
