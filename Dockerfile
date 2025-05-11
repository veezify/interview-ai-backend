# Build stage
FROM golang:1.23-alpine AS builder
# Build stage

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# List files to debug
RUN ls -la
RUN ls -la cmd/api

# Build the application with more verbose output
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o interview-ai-backend ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/interview-ai-backend .
# Copy .env file if needed (optional)
COPY --from=builder /app/.env .

# Expose the application port
EXPOSE 8080

# Command to run
CMD ["./interview-ai-backend"]