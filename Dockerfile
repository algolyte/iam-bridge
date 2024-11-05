# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

# Install Wire
RUN go install github.com/google/wire/cmd/wire@latest

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Generate wire_gen.go
RUN wire ./internal/di

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .
COPY --from=builder /app/config ./config

# Install CA certificates
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -g '' appuser
USER appuser

CMD ["./main"]
