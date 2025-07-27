# Stage 1: Builder
FROM golang:1.23.4-alpine3.19 AS builder

WORKDIR /app

# Install git (needed for go mod sometimes)
RUN apk add --no-cache git

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o system_service ./cmd/server/

# Stage 2: Runner
FROM alpine:3.19.1

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app .

# Expose port if needed (e.g., 8080)
EXPOSE 8080

# Run the binary
CMD ["./system_service"]
