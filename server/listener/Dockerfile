FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o listener ./cmd/api

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/listener .

CMD ["./listener"]