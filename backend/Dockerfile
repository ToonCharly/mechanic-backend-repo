# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and source files
COPY go.mod ./
COPY . .

# Download dependencies and generate go.sum
RUN go mod download
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Final stage - Production hardened
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -D -u 1000 -G appgroup appuser

# Set working directory
WORKDIR /home/appuser

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy .env file (optional, better to use env vars in docker-compose)
COPY .env.example .env 

# Change ownership to appuser
RUN chown -R appuser:appgroup /home/appuser

# Switch to non-root user
USER appuser

EXPOSE 8080

CMD ["./main"]
