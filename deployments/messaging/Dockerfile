FROM golang:1.24.1-alpine AS builder

RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build  -o messaging-service ./cmd/messaging

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/messaging-service .

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

CMD ["./messaging-service"]
