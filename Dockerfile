# syntax=docker/dockerfile:1

FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor

COPY . .

RUN apk add --no-cache openssh

RUN go build -o main cmd/api/application.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

EXPOSE 8080

CMD ["./main"]
