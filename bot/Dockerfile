# syntax=docker/dockerfile:1

# --- Build stage ---
FROM golang:1.24 AS builder
WORKDIR /app

# Download deps first (better layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source
COPY . .

# Build a static binary
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bot ./

# --- Runtime stage ---
FROM alpine:3.20

# Certificates for HTTPS (Telegram API) + non-root user
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app

WORKDIR /app
COPY --from=builder /app/bot /app/bot

USER app

CMD ["./bot"]
