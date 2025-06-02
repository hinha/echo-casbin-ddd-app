# Use the official Golang image as the base image
FROM golang:1.20-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/main.go

# Use a minimal alpine image for the final image
FROM alpine:3.17

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Copy web files
COPY web/ /app/web/

# Copy casbin model file
COPY casbin/ /app/casbin/

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./app"]