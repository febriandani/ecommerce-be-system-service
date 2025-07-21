# Use an official Golang image as the base image
FROM alpine:1.23 as builder

WORKDIR /app

# Copy Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o user-service .

# Use a minimal base image for the final container
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/system_service .
CMD ["./system_service"]
