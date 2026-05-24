# syntax=docker/dockerfile:1
# Build context: backend/
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /workspace

COPY go.work go.work.sum ./
COPY go.mod go.sum ./
COPY services/api-gateway/go.mod      services/api-gateway/go.mod
COPY services/auth-service/go.mod     services/auth-service/go.mod
COPY services/trip-service/go.mod     services/trip-service/go.mod
COPY services/user-service/go.mod     services/user-service/go.mod
COPY services/collaboration-service/go.mod services/collaboration-service/go.mod
COPY services/notification-service/go.mod  services/notification-service/go.mod
COPY services/search-service/go.mod        services/search-service/go.mod

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app ./services/auth-service/cmd/...

FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app /app/auth-service
COPY --chown=appuser:appgroup services/auth-service/config ./config

USER appuser
EXPOSE 8081
ENTRYPOINT ["/app/auth-service"]
