# Use Golang image for the backend
FROM golang:1.23-alpine as builder

# Set the current working directory inside the container
WORKDIR /app

# Copy the Go files
COPY . .

# Install dependencies and build the Go app
RUN go mod tidy
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
