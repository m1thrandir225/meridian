FROM golang:1.24.1-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /app/messaging-server ./cmd/messaging

FROM alpine:latest

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 8081

CMD ["/app/messaging-server"]
