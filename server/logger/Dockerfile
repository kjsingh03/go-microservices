FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o logger ./cmd/api

FROM alpine:latest

COPY --from=builder /app/logger .

EXPOSE 8085

CMD ["./logger"]