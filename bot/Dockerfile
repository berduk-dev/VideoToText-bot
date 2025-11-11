# syntax=docker/dockerfile:1
FROM golang:1.24 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./bot ./bot

# Собираем бинарь в корень /app, а не в папку /app/bot
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/bot_bin ./bot/cmd/bot

# --- Runtime stage ---
FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app

WORKDIR /app

# Копируем бинарь, не папку
COPY --from=builder /app/bot_bin /app/bot
USER app

CMD ["./bot"]
