# Этап 1: Сборка бинарного файла
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o Sender ./cmd/sender/main.go

# Этап 2: Создание минимального образа
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/Sender .

CMD ["sh", "-c", "sleep 20 && ./Sender"]
