FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o broker-service ./cmd/api

FROM alpine:latest  

WORKDIR /root/

COPY --from=builder /app/broker-service .

EXPOSE 8084

CMD ["./broker-service"]