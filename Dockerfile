FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd

FROM alpine:3.20

RUN apk add --no-cache ca-certificates && adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /app/bin/server /app/server

USER appuser

EXPOSE 8080 9090

ENTRYPOINT ["/app/server"]
