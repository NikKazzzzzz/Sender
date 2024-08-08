FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY cmd/ .
COPY config/sender.yaml ./config/sender.yaml

RUN go mod tidy
RUN go build -o sender

FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/sender .
COPY --from=builder /app/config/sender.yaml ./config/sender.yaml

CMD ["./sender"]
