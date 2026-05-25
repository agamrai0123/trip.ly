# ctx:auth-service | 2026-05-24 | HTTP:8081 gRPC:9081
## state
todo: Dockerfile — multi-stage(golang:1.22-alpine→gcr.io/distroless/static-debian12), CGO_ENABLED=0 GOOS=linux -ldflags="-w -s", non-root, EXPOSE 8081 9081, HEALTHCHECK /healthz, LABEL service=auth-service
planned: -
done: scaffold(cmd+internal/*)@2026-05-24 | env-var-fix(DB_*+OAuth+KAFKA_BROKERS flat-env overrides)@2026-05-24
errors: -
_update when: OAuth providers/JWT config/endpoints/users+refresh_tokens schema/Kafka events/gRPC RPCs change_

## files
cmd/main.go: bootstrap — DB pool,Kafka producer,JWT manager,OAuth2 configs,gRPC server,gin,shutdown(SIGTERM/SIGINT 10s)
config/auth-service-config.toml: TOML defaults
internal/config.go: Viper; OAuth IDs/secrets,Kafka brokers,JWT key paths,DB DSN
internal/service.go: OAuthService(InitiateLogin,HandleCallback,RefreshToken,Logout,GetMe); ValidateToken(gRPC impl); user upsert
internal/handlers.go: LoginHandler,CallbackHandler,RefreshHandler,LogoutHandler,MeHandler
internal/routes.go: RegisterRoutes; /auth/* group
internal/models.go: User,RefreshToken,OAuthConfig,UserClaims structs
internal/database.go: pgx UserRepo+RefreshTokenRepo; all parameterised queries
internal/metrics.go: oauth_login_total,token_refresh_total,token_validate_total,db_query_duration,goroutine_count; /metrics
internal/logger.go: zerolog+lumberjack | internal/errors.go: AppError helpers
certs/: RS256 private+public key PEM files (dev only)

## routes
GET /healthz,/metrics [no-auth]
GET /auth/:provider/login [no-auth] — redirect to OAuth provider (state+PKCE)
GET /auth/:provider/callback [no-auth] — exchange code → JWT + httpOnly refresh cookie
POST /auth/refresh [httpOnly-cookie] — rotate refresh token → new JWT
POST /auth/logout [JWT] — revoke refresh token
GET /auth/me [JWT] — return user from JWT claims

## grpc-in: AuthService.ValidateToken(token) → {user_id,email,name,avatar_url}; called-by: api-gateway (every JWT req)
## kafka-out: auth-events:user.login {user_id,email,provider,timestamp} on OAuth callback success

## env: PORT=8081 GRPC_PORT=9081 DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD KAFKA_BROKERS JWT_PRIVATE_KEY(b64) JWT_PUBLIC_KEY(b64) JWT_EXPIRY=15m GOOGLE_CLIENT_ID GOOGLE_CLIENT_SECRET GOOGLE_REDIRECT_URL GITHUB_CLIENT_ID GITHUB_CLIENT_SECRET GITHUB_REDIRECT_URL FRONTEND_URL
## pkg: config database errors jwt kafka logger middleware response
## db: users(upsert on OAuth callback,SELECT for GetMe) | refresh_tokens(INSERT on issue,SELECT+UPDATE revoked_at on refresh/logout)
## security: refresh-token stored as SHA-256(token) never plaintext; rotated on every use; httpOnly;Secure;SameSite=Strict cookie; access token in-memory only; RS256 private key only in auth-service



