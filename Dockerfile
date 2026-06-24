# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
# CGO_ENABLED=0 is used to build a static binary, perfect for Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app-binary ./cmd/api

# Final stage
FROM alpine:latest

WORKDIR /app

# Install tzdata if your app needs timezone info
RUN apk --no-cache add tzdata

# Copy the binary from the builder stage
COPY --from=builder /app-binary /app/api

# Expose the application port
EXPOSE 8080

# Command to run the executable
CMD ["/app/api"]
