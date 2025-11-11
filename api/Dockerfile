FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# копируем весь проект (чтобы подтянуть api/)
COPY . .

# теперь путь правильный
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o app ./api/internal/cmd

CMD ["./app"]
