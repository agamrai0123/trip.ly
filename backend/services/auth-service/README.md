# auth-service

OAuth 2.0 PKCE authentication service for WanderPlan. Issues RS256 JWTs and manages refresh token rotation.

## Purpose

- Google and GitHub OAuth 2.0 PKCE login flows
- Issues short-lived (15-min) RS256 JWTs stored in-memory on the client
- Opaque refresh tokens stored as SHA-256 hashes in PostgreSQL, rotated on every use
- Exposes gRPC `ValidateToken` RPC consumed by api-gateway for every authenticated request
- Publishes `user.login` events to Kafka topic `auth-events`

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/auth/:provider/login` | None | Redirect to OAuth provider (`google` or `github`) |
| `GET` | `/auth/:provider/callback` | None | Exchange code → JWT + httpOnly refresh cookie |
| `POST` | `/auth/refresh` | httpOnly cookie | Rotate refresh token → new JWT |
| `POST` | `/auth/logout` | JWT | Revoke refresh token |
| `GET` | `/auth/me` | JWT | Return user from JWT claims |

## gRPC

- **Port:** `9081`
- **Service:** `AuthService`
- **RPC:** `ValidateToken(token) → {user_id, email, name, avatar_url}`
- Called by api-gateway on every JWT-protected request

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP port (default: `8081`) |
| `GRPC_PORT` | gRPC port (default: `9081`) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_NAME` | Database name |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_SSL_MODE` | `disable` / `require` |
| `KAFKA_BROKERS` | Comma-separated Kafka broker addresses |
| `JWT_PRIVATE_KEY` | Base64-encoded RSA private key PEM |
| `JWT_PUBLIC_KEY` | Base64-encoded RSA public key PEM |
| `JWT_EXPIRY` | Token expiry duration (default: `15m`) |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID |
| `GOOGLE_CLIENT_SECRET` | Google OAuth client secret |
| `GOOGLE_REDIRECT_URL` | Google OAuth redirect URL |
| `GITHUB_CLIENT_ID` | GitHub OAuth client ID |
| `GITHUB_CLIENT_SECRET` | GitHub OAuth client secret |
| `GITHUB_REDIRECT_URL` | GitHub OAuth redirect URL |
| `FRONTEND_URL` | Frontend URL for post-login redirect |
| `RATE_LIMIT_RPS` | Rate limit — requests per second |
| `RATE_LIMIT_BURST` | Rate limit — burst size |

## How to Run Locally

```bash
cp .env.example .env   # edit OAuth credentials and DB values
cd backend/services/auth-service
air                    # hot-reload via .air.toml
```

## Generate RSA Keys (dev)

```bash
openssl genrsa -out private.pem 2048
openssl rsa -in private.pem -pubout -out public.pem
# Then base64-encode each for env vars:
base64 -w0 private.pem   # → JWT_PRIVATE_KEY
base64 -w0 public.pem    # → JWT_PUBLIC_KEY
```

## How to Run Tests

```bash
cd backend/services/auth-service
go test ./... -v -count=1
```

Integration tests require Docker (PostgreSQL via testcontainers-go):

```bash
go test ./... -v -count=1 -tags integration
```

## Observability

- `GET /metrics` — `oauth_login_total`, `token_refresh_total`, `token_validate_total`, `db_query_duration_seconds`, `go_goroutines`
- Structured JSON logs via `zerolog` + `lumberjack`

## Security Notes

- Refresh tokens are **never** stored in plaintext — SHA-256 hash only
- Access tokens are kept in-memory on the client; never in localStorage
- Refresh token cookie: `httpOnly`, `Secure`, `SameSite=Strict`
- RS256 private key never leaves this service
