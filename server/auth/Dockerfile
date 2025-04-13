FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o auth-service ./cmd/api

FROM alpine:latest

WORKDIR /root

COPY --from=builder /app/auth-service .

EXPOSE 8081

CMD [ "./auth-service" ]