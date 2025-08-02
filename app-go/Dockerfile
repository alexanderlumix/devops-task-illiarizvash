# syntax=docker/dockerfile:1

# Multi-stage build for optimized production image
# Stage 1: Build stage
FROM golang:1.21-alpine AS build

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies (cached layer)
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimizations
# - CGO_ENABLED=0: Static linking for smaller image
# - GOOS=linux: Target Linux platform
# - -ldflags: Strip debug info and optimize binary size
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o read_products read_products.go

# Stage 2: Production stage
# Use minimal alpine image for security and size
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /app

# Copy binary from build stage
COPY --from=build /app/read_products .

# Change ownership to non-root user
RUN chown -R appuser:appgroup /app

# Switch to non-root user
USER appuser

# Expose port (if needed for health checks)
EXPOSE 8080

# Health check to ensure application is running
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ps aux | grep read_products || exit 1

# Set environment variables
ENV GIN_MODE=release

# Run the application
CMD ["./read_products"]
