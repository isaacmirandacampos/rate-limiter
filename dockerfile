FROM golang:1.23.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o rate-limiter-api cmd/rate-limiter-api/main.go
COPY /.env.dockerfile .env

FROM alpine:latest

# Copy binary from builder
COPY --from=builder /app/rate-limiter-api /rate-limiter-api
COPY --from=builder /app/.env /.env

# Run
CMD ["/rate-limiter-api"]