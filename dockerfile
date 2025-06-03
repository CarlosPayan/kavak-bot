FROM golang:1.24-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o kavak-bot cmd/bot/main.go

FROM alpine:latest
WORKDIR /root

COPY --from=builder /app/kavak-bot .
COPY --from=builder /app/configs/config.yaml configs/config.yaml
COPY --from=builder /app/data/catalog.csv data/catalog.csv

EXPOSE 8080

ENTRYPOINT ["./kavak-bot"]
