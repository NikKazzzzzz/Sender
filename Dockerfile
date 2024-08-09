# Dockerfile for Sender

# Use a Golang base image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o sender cmd/sender/main.go

# Command to run the application
CMD ["./sender"]
