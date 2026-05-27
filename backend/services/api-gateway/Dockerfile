# syntax=docker/dockerfile:1
# Build context: backend/
FROM golang:1.26-alpine AS builder
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /workspace

# Copy workspace files for dependency resolution
COPY go.work go.work.sum ./
COPY go.mod go.sum ./

# Copy all service go.mod files for workspace modules
COPY services/api-gateway/go.mod      services/api-gateway/go.mod
COPY services/auth-service/go.mod     services/auth-service/go.mod
COPY services/trip-service/go.mod     services/trip-service/go.mod
COPY services/user-service/go.mod     services/user-service/go.mod
COPY services/collaboration-service/go.mod services/collaboration-service/go.mod
COPY services/notification-service/go.mod  services/notification-service/go.mod
COPY services/search-service/go.mod        services/search-service/go.mod

# Copy all source (workspace-local packages cannot be fetched remotely)
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app ./services/api-gateway/cmd/...

# ── Final stage ────────────────────────────────────────────────────────────────
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S appgroup && adduser -S appuser -G appgroup
LABEL service="api-gateway" version="1.0.0" maintainer="WanderPlan"
WORKDIR /app
COPY --from=builder /app /app/api-gateway
COPY --chown=appuser:appgroup services/api-gateway/config ./config

USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/app/api-gateway"]
