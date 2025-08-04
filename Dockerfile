# Use a multi-stage build for better efficiency
FROM golang:1.24.2-bullseye AS builder

# Install Erlang and build dependencies
RUN apt-get update && apt-get install -y \
    erlang-dev \
    erlang-ei \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Set CGO flags for Erlang interface
ENV CGO_CFLAGS="-I/usr/lib/erlang/lib/erl_interface-*/include -Wall -g"
ENV CGO_LDFLAGS="-L/usr/lib/erlang/lib/erl_interface-*/lib -lei -lpthread"

# Build the application
RUN CGO_ENABLED=1 go build -o gotp .

# Runtime stage
FROM debian:bullseye-slim

# Install runtime Erlang dependencies
RUN apt-get update && apt-get install -y \
    erlang-base \
    erlang-ei \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN useradd -r -u 1001 -g root gotp

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/gotp .

# Change ownership to non-root user
RUN chown gotp:root /app/gotp

# Switch to non-root user
USER gotp

# Expose port (if needed for the application)
EXPOSE 9999

# Run the application
CMD ["./gotp"]