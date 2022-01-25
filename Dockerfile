# syntax=docker/dockerfile:1
# Build Stage
FROM golang:1.17.3-alpine3.14 AS builder

RUN apk update
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd


# Run Stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates
RUN apk add bash


WORKDIR /app

COPY --from=builder /app/main .

COPY --from=builder /app/.env .

COPY wait-for-it.sh .

RUN chmod +x wait-for-it.sh


CMD ["./main"]