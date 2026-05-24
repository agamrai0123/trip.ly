---
applyTo: "**/Dockerfile"
---

# Dockerfile Rules

## Structure
- Always use multi-stage builds: a `builder` stage and a minimal `final` stage.
- Builder stage: `FROM golang:1.22-alpine AS builder`. Final stage: `FROM gcr.io/distroless/static-debian12` or `FROM alpine:3.20`.
- Copy only the compiled binary into the final stage — never copy the full Go source.

## Go build flags
- Build with: `CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/service ./cmd/...`
- The `-w -s` flags strip debug info and reduce binary size.

## Security
- Run the process as a non-root user in the final stage: `USER nonroot:nonroot` (distroless) or create a dedicated user with `adduser`.
- Never copy `.env` files into the image. All config comes from environment variables at runtime.
- Pin base image versions — never use `latest`.

## Metadata
- Add `LABEL` with `service`, `version`, and `maintainer`.
- Expose only the service's HTTP and gRPC ports with `EXPOSE`.

## Health check
- Add a `HEALTHCHECK` that calls `GET /healthz` using `wget` or `curl` (alpine) or the built-in binary (distroless).
