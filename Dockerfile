# Build stage
FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the CLI application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o wedl wedl.go

# Build the API server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api-server api-server.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy both binaries from builder stage
COPY --from=builder /app/wedl .
COPY --from=builder /app/api-server .

# Create downloads directory
RUN mkdir -p /downloads

# Set the binaries as executable
RUN chmod +x ./wedl ./api-server

# Expose API port
EXPOSE 8080

# Default command runs API server
ENTRYPOINT ["/root/api-server"]
CMD ["-port", "8080"]