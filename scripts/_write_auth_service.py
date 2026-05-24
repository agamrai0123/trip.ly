import os

base = r'D:\Learn\trip.ly\backend\services\auth-service\internal'

database_go = '''package internal

import (
\t"context"
\t"crypto/sha256"
\t"encoding/hex"
\t"errors"
\t"time"

\t"github.com/google/uuid"
\t"github.com/jackc/pgx/v5"
\t"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a row is not found.
var ErrNotFound = errors.New("not found")

// UserRepo handles user persistence.
type UserRepo struct{ pool *pgxpool.Pool }

func NewUserRepo(pool *pgxpool.Pool) *UserRepo { return &UserRepo{pool: pool} }

// Upsert inserts or updates a user row and returns the current record.
func (r *UserRepo) Upsert(ctx context.Context, u *User) (*User, error) {
\trow := r.pool.QueryRow(ctx, `
\t\tINSERT INTO wanderplan.users (id, email, name, avatar_url, provider, created_at, updated_at)
\t\tVALUES ($1, $2, $3, $4, $5, NOW(), NOW())
\t\tON CONFLICT (email) DO UPDATE
\t\t\tSET name       = EXCLUDED.name,
\t\t\t    avatar_url = EXCLUDED.avatar_url,
\t\t\t    provider   = EXCLUDED.provider,
\t\t\t    updated_at = NOW()
\t\tRETURNING id, email, name, avatar_url, provider, created_at, updated_at`,
\t\tu.ID, u.Email, u.Name, u.AvatarURL, u.Provider,
\t)
\tout := &User{}
\terr := row.Scan(&out.ID, &out.Email, &out.Name, &out.AvatarURL, &out.Provider, &out.CreatedAt, &out.UpdatedAt)
\tif err != nil {
\t\treturn nil, err
\t}
\treturn out, nil
}

// GetByID fetches a user by their UUID.
func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
\trow := r.pool.QueryRow(ctx, `
\t\tSELECT id, email, name, avatar_url, provider, created_at, updated_at
\t\tFROM wanderplan.users WHERE id = $1`, id)
\tu := &User{}
\terr := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.Provider, &u.CreatedAt, &u.UpdatedAt)
\tif errors.Is(err, pgx.ErrNoRows) {
\t\treturn nil, ErrNotFound
\t}
\tif err != nil {
\t\treturn nil, err
\t}
\treturn u, nil
}

// RefreshTokenRepo handles refresh token persistence.
type RefreshTokenRepo struct{ pool *pgxpool.Pool }

func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
\treturn &RefreshTokenRepo{pool: pool}
}

func hashToken(token string) string {
\th := sha256.Sum256([]byte(token))
\treturn hex.EncodeToString(h[:])
}

// Create stores a new refresh token and returns the raw opaque token string.
func (r *RefreshTokenRepo) Create(ctx context.Context, userID string, ttl time.Duration) (string, error) {
\traw := uuid.New().String() + uuid.New().String()
\ttokenHash := hashToken(raw)
\t_, err := r.pool.Exec(ctx, `
\t\tINSERT INTO wanderplan.refresh_tokens (id, user_id, token_hash, expires_at, created_at)
\t\tVALUES ($1, $2, $3, $4, NOW())`,
\t\tuuid.New().String(), userID, tokenHash, time.Now().Add(ttl),
\t)
\tif err != nil {
\t\treturn "", err
\t}
\treturn raw, nil
}

// Rotate validates a token, deletes it (rotation), and creates a new one.
func (r *RefreshTokenRepo) Rotate(ctx context.Context, rawToken string, ttl time.Duration) (string, string, error) {
\ttokenHash := hashToken(rawToken)
\tvar userID string
\tvar expiresAt time.Time
\terr := r.pool.QueryRow(ctx, `
\t\tSELECT user_id, expires_at FROM wanderplan.refresh_tokens
\t\tWHERE token_hash = $1`, tokenHash,
\t).Scan(&userID, &expiresAt)
\tif errors.Is(err, pgx.ErrNoRows) {
\t\treturn "", "", ErrNotFound
\t}
\tif err != nil {
\t\treturn "", "", err
\t}
\tif time.Now().After(expiresAt) {
\t\t_, _ = r.pool.Exec(ctx, `DELETE FROM wanderplan.refresh_tokens WHERE token_hash = $1`, tokenHash)
\t\treturn "", "", ErrNotFound
\t}
\tif _, err = r.pool.Exec(ctx, `DELETE FROM wanderplan.refresh_tokens WHERE token_hash = $1`, tokenHash); err != nil {
\t\treturn "", "", err
\t}
\tnewRaw, err := r.Create(ctx, userID, ttl)
\tif err != nil {
\t\treturn "", "", err
\t}
\treturn userID, newRaw, nil
}

// DeleteByUserID revokes all refresh tokens for a user (logout).
func (r *RefreshTokenRepo) DeleteByUserID(ctx context.Context, userID string) error {
\t_, err := r.pool.Exec(ctx, `DELETE FROM wanderplan.refresh_tokens WHERE user_id = $1`, userID)
\treturn err
}
'''

service_go = '''package internal

import (
\t"context"
\t"crypto/rand"
\t"encoding/base64"
\t"encoding/json"
\t"fmt"
\t"net/http"
\t"time"

\t"github.com/rs/zerolog/log"
\t"golang.org/x/oauth2"
\t"golang.org/x/oauth2/github"
\t"golang.org/x/oauth2/google"
\t"wanderplan/pkg/jwt"
\t"wanderplan/pkg/kafka"
)

const (
\trefreshTokenTTL = 30 * 24 * time.Hour
\taccessTokenTTL  = 15 * time.Minute
)

// AuthService orchestrates OAuth PKCE, JWT issuance, and Kafka events.
type AuthService struct {
\tcfg           *Config
\tusers         *UserRepo
\trefreshTokens *RefreshTokenRepo
\tjwtMgr        *jwt.Manager
\tproducer      *kafka.Producer
\tgoogleCfg     *oauth2.Config
\tgithubCfg     *oauth2.Config
}

// NewAuthService constructs the service with all dependencies.
func NewAuthService(
\tcfg *Config,
\tusers *UserRepo,
\trt *RefreshTokenRepo,
\tjwtMgr *jwt.Manager,
\tproducer *kafka.Producer,
) *AuthService {
\tgCfg := &oauth2.Config{
\t\tClientID:     cfg.OAuth.Google.ClientID,
\t\tClientSecret: cfg.OAuth.Google.ClientSecret,
\t\tRedirectURL:  cfg.OAuth.Google.RedirectURL,
\t\tScopes:       []string{"openid", "email", "profile"},
\t\tEndpoint:     google.Endpoint,
\t}
\tghCfg := &oauth2.Config{
\t\tClientID:     cfg.OAuth.GitHub.ClientID,
\t\tClientSecret: cfg.OAuth.GitHub.ClientSecret,
\t\tRedirectURL:  cfg.OAuth.GitHub.RedirectURL,
\t\tScopes:       []string{"read:user", "user:email"},
\t\tEndpoint:     github.Endpoint,
\t}
\treturn &AuthService{
\t\tcfg:           cfg,
\t\tusers:         users,
\t\trefreshTokens: rt,
\t\tjwtMgr:        jwtMgr,
\t\tproducer:      producer,
\t\tgoogleCfg:     gCfg,
\t\tgithubCfg:     ghCfg,
\t}
}

// BuildAuthURL generates the OAuth2 PKCE authorization URL.
// Returns (url, state, codeVerifier, error).
func (s *AuthService) BuildAuthURL(provider string) (string, string, string, error) {
\tstate, err := randomBase64(32)
\tif err != nil {
\t\treturn "", "", "", fmt.Errorf("generate state: %w", err)
\t}
\tverifier, err := randomBase64(32)
\tif err != nil {
\t\treturn "", "", "", fmt.Errorf("generate verifier: %w", err)
\t}
\tcfg, err := s.providerConfig(provider)
\tif err != nil {
\t\treturn "", "", "", err
\t}
\turl := cfg.AuthCodeURL(state,
\t\toauth2.AccessTypeOffline,
\t\toauth2.S256ChallengeOption(verifier),
\t)
\treturn url, state, verifier, nil
}

// ExchangeCode swaps the authorization code for user info and issues WanderPlan tokens.
func (s *AuthService) ExchangeCode(
\tctx context.Context,
\tprovider, code, verifier string,
) (*LoginResponse, string, error) {
\tcfg, err := s.providerConfig(provider)
\tif err != nil {
\t\treturn nil, "", err
\t}
\ttoken, err := cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
\tif err != nil {
\t\treturn nil, "", fmt.Errorf("oauth exchange: %w", err)
\t}
\tuserInfo, err := s.fetchUserInfo(ctx, provider, token)
\tif err != nil {
\t\treturn nil, "", err
\t}
\tuser, err := s.users.Upsert(ctx, userInfo)
\tif err != nil {
\t\treturn nil, "", fmt.Errorf("upsert user: %w", err)
\t}
\taccessToken, err := s.jwtMgr.Sign(user.ID, user.Email, user.Name, user.AvatarURL)
\tif err != nil {
\t\treturn nil, "", fmt.Errorf("sign jwt: %w", err)
\t}
\trawRefresh, err := s.refreshTokens.Create(ctx, user.ID, refreshTokenTTL)
\tif err != nil {
\t\treturn nil, "", fmt.Errorf("create refresh token: %w", err)
\t}
\tpayload, _ := json.Marshal(map[string]string{
\t\t"user_id": user.ID,
\t\t"email":   user.Email,
\t\t"event":   "user.login",
\t})
\tif err := s.producer.Publish(ctx, kafka.TopicAuthEvents, "user.login", payload); err != nil {
\t\tlog.Warn().Err(err).Msg("failed to publish auth event")
\t}
\treturn &LoginResponse{AccessToken: accessToken, User: user}, rawRefresh, nil
}

// Refresh validates a refresh token, rotates it, and issues a new access token.
func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*RefreshResponse, string, error) {
\tuserID, newRaw, err := s.refreshTokens.Rotate(ctx, rawToken, refreshTokenTTL)
\tif err != nil {
\t\treturn nil, "", err
\t}
\tuser, err := s.users.GetByID(ctx, userID)
\tif err != nil {
\t\treturn nil, "", err
\t}
\taccessToken, err := s.jwtMgr.Sign(user.ID, user.Email, user.Name, user.AvatarURL)
\tif err != nil {
\t\treturn nil, "", err
\t}
\treturn &RefreshResponse{AccessToken: accessToken}, newRaw, nil
}

// Logout revokes all refresh tokens for the user.
func (s *AuthService) Logout(ctx context.Context, userID string) error {
\treturn s.refreshTokens.DeleteByUserID(ctx, userID)
}

// ValidateToken validates a JWT and returns the claims.
func (s *AuthService) ValidateToken(tokenStr string) (*jwt.Claims, error) {
\treturn s.jwtMgr.Parse(tokenStr)
}

func (s *AuthService) providerConfig(provider string) (*oauth2.Config, error) {
\tswitch provider {
\tcase "google":
\t\treturn s.googleCfg, nil
\tcase "github":
\t\treturn s.githubCfg, nil
\tdefault:
\t\treturn nil, fmt.Errorf("unknown provider: %s", provider)
\t}
}

type googleUserInfo struct {
\tSub     string `json:"sub"`
\tEmail   string `json:"email"`
\tName    string `json:"name"`
\tPicture string `json:"picture"`
}

type githubUserInfo struct {
\tLogin     string `json:"login"`
\tEmail     string `json:"email"`
\tAvatarURL string `json:"avatar_url"`
\tName      string `json:"name"`
}

func (s *AuthService) fetchUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*User, error) {
\tclient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
\tswitch provider {
\tcase "google":
\t\tresp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tdefer resp.Body.Close()
\t\tif resp.StatusCode != http.StatusOK {
\t\t\treturn nil, fmt.Errorf("google userinfo status %d", resp.StatusCode)
\t\t}
\t\tvar info googleUserInfo
\t\tif err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
\t\t\treturn nil, err
\t\t}
\t\treturn &User{ID: info.Sub, Email: info.Email, Name: info.Name, AvatarURL: info.Picture, Provider: "google"}, nil
\tcase "github":
\t\tresp, err := client.Get("https://api.github.com/user")
\t\tif err != nil {
\t\t\treturn nil, err
\t\t}
\t\tdefer resp.Body.Close()
\t\tif resp.StatusCode != http.StatusOK {
\t\t\treturn nil, fmt.Errorf("github user status %d", resp.StatusCode)
\t\t}
\t\tvar info githubUserInfo
\t\tif err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
\t\t\treturn nil, err
\t\t}
\t\tif info.Email == "" {
\t\t\tinfo.Email, _ = fetchGitHubPrimaryEmail(client)
\t\t}
\t\tname := info.Name
\t\tif name == "" {
\t\t\tname = info.Login
\t\t}
\t\treturn &User{ID: fmt.Sprintf("gh_%s", info.Login), Email: info.Email, Name: name, AvatarURL: info.AvatarURL, Provider: "github"}, nil
\t}
\treturn nil, fmt.Errorf("unknown provider: %s", provider)
}

type ghEmail struct {
\tEmail   string `json:"email"`
\tPrimary bool   `json:"primary"`
}

func fetchGitHubPrimaryEmail(client *http.Client) (string, error) {
\tresp, err := client.Get("https://api.github.com/user/emails")
\tif err != nil {
\t\treturn "", err
\t}
\tdefer resp.Body.Close()
\tvar emails []ghEmail
\tif err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
\t\treturn "", err
\t}
\tfor _, e := range emails {
\t\tif e.Primary {
\t\t\treturn e.Email, nil
\t\t}
\t}
\treturn "", nil
}

func randomBase64(n int) (string, error) {
\tb := make([]byte, n)
\tif _, err := rand.Read(b); err != nil {
\t\treturn "", err
\t}
\treturn base64.RawURLEncoding.EncodeToString(b), nil
}
'''

handlers_go = '''package internal

import (
\t"net/http"
\t"time"

\t"github.com/gin-gonic/gin"
\t"wanderplan/pkg/errors"
\t"wanderplan/pkg/middleware"
\t"wanderplan/pkg/response"
)

// Handlers wires HTTP requests to AuthService methods.
type Handlers struct {
\tsvc *AuthService
\tcfg *Config
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *AuthService, cfg *Config) *Handlers {
\treturn &Handlers{svc: svc, cfg: cfg}
}

// Login redirects the user to the OAuth provider.
// GET /auth/:provider/login
func (h *Handlers) Login(c *gin.Context) {
\tprovider := c.Param("provider")
\tauthURL, state, verifier, err := h.svc.BuildAuthURL(provider)
\tif err != nil {
\t\tresponse.Err(c, errors.BadRequest(err.Error()))
\t\treturn
\t}
\tmaxAge := h.cfg.Cookie.StateMaxAgeSecs
\tif maxAge == 0 {
\t\tmaxAge = 300
\t}
\tc.SetSameSite(http.SameSiteLaxMode)
\tc.SetCookie("oauth_state", state, maxAge, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
\tc.SetCookie("pkce_verifier", verifier, maxAge, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
\tc.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles the OAuth provider redirect.
// GET /auth/:provider/callback
func (h *Handlers) Callback(c *gin.Context) {
\tprovider := c.Param("provider")
\tcookieState, err := c.Cookie("oauth_state")
\tif err != nil || cookieState != c.Query("state") {
\t\tresponse.Err(c, errors.BadRequest("invalid oauth state"))
\t\treturn
\t}
\tverifier, err := c.Cookie("pkce_verifier")
\tif err != nil {
\t\tresponse.Err(c, errors.BadRequest("missing pkce verifier"))
\t\treturn
\t}
\tc.SetCookie("oauth_state", "", -1, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
\tc.SetCookie("pkce_verifier", "", -1, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)

\tloginResp, rawRefresh, err := h.svc.ExchangeCode(c.Request.Context(), provider, c.Query("code"), verifier)
\tif err != nil {
\t\tresponse.Err(c, errors.Internal(err.Error()))
\t\treturn
\t}
\th.setRefreshCookie(c, rawRefresh)
\tresponse.OK(c, loginResp)
}

// Refresh issues a new access token using the refresh_token cookie.
// POST /auth/refresh
func (h *Handlers) Refresh(c *gin.Context) {
\trawToken, err := c.Cookie("refresh_token")
\tif err != nil {
\t\tresponse.Err(c, errors.Unauthorized("missing refresh token"))
\t\treturn
\t}
\trefreshResp, newRaw, err := h.svc.Refresh(c.Request.Context(), rawToken)
\tif err != nil {
\t\th.clearRefreshCookie(c)
\t\tresponse.Err(c, errors.Unauthorized("invalid or expired refresh token"))
\t\treturn
\t}
\th.setRefreshCookie(c, newRaw)
\tresponse.OK(c, refreshResp)
}

// Logout revokes the refresh token and clears the cookie.
// POST /auth/logout
func (h *Handlers) Logout(c *gin.Context) {
\tuserID := middleware.UserID(c)
\tif userID != "" {
\t\t_ = h.svc.Logout(c.Request.Context(), userID)
\t}
\th.clearRefreshCookie(c)
\tresponse.NoContent(c)
}

// Me returns the authenticated user\'s claims.
// GET /auth/me
func (h *Handlers) Me(c *gin.Context) {
\tuserID := middleware.UserID(c)
\temail, _ := c.Get("email")
\tname, _ := c.Get("name")
\tavatar, _ := c.Get("avatar_url")
\tresponse.OK(c, gin.H{
\t\t"user_id":    userID,
\t\t"email":      email,
\t\t"name":       name,
\t\t"avatar_url": avatar,
\t})
}

func (h *Handlers) setRefreshCookie(c *gin.Context, raw string) {
\tconst maxAge = int(30 * 24 * time.Hour / time.Second)
\tc.SetSameSite(http.SameSiteStrictMode)
\tc.SetCookie("refresh_token", raw, maxAge, "/auth", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
}

func (h *Handlers) clearRefreshCookie(c *gin.Context) {
\tc.SetCookie("refresh_token", "", -1, "/auth", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
}
'''

routes_go = '''package internal

import (
\t"github.com/gin-gonic/gin"
\t"github.com/prometheus/client_golang/prometheus"
\t"wanderplan/pkg/jwt"
\t"wanderplan/pkg/middleware"
)

// RegisterRoutes mounts all auth-service routes on the engine.
func RegisterRoutes(r *gin.Engine, h *Handlers, jwtMgr *jwt.Manager, reg *prometheus.Registry, cfg *Config) {
\tmetricsMW, metricsHandler := middleware.Metrics("auth_service", reg)

\tr.Use(
\t\tmiddleware.Recovery(),
\t\tmiddleware.RequestID(),
\t\tmiddleware.Logger(),
\t\tmiddleware.CORS(cfg.CORS.AllowedOrigins),
\t\tmiddleware.RateLimit(cfg.RateLimit.RPS, cfg.RateLimit.Burst),
\t\tmetricsMW,
\t)

\tr.GET("/metrics", gin.WrapH(metricsHandler))
\tr.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

\tauth := r.Group("/auth")
\t{
\t\tauth.GET("/:provider/login", h.Login)
\t\tauth.GET("/:provider/callback", h.Callback)
\t\tauth.POST("/refresh", h.Refresh)
\t\tauth.POST("/logout", middleware.Auth(jwtMgr), h.Logout)
\t\tauth.GET("/me", middleware.Auth(jwtMgr), h.Me)
\t}
}
'''

errors_go = '''package internal

// Re-export pkg/errors constructors for convenience.

import pkgerr "wanderplan/pkg/errors"

var (
\tBadRequest   = pkgerr.BadRequest
\tUnauthorized = pkgerr.Unauthorized
\tInternal     = pkgerr.Internal
)
'''

logger_go = '''package internal

import (
\t"os"

\t"github.com/rs/zerolog"
\t"github.com/rs/zerolog/log"
\tpkglogger "wanderplan/pkg/logger"
)

// InitLogger initialises the global zerolog logger for the auth-service.
func InitLogger(cfg LoggingCfg, service string) zerolog.Logger {
\tlevel := zerolog.Level(cfg.Level)
\tif level < zerolog.TraceLevel || level > zerolog.Disabled {
\t\tlevel = zerolog.InfoLevel
\t}
\tmaxSizeMB := cfg.MaxSizeMB
\tif maxSizeMB <= 0 {
\t\tmaxSizeMB = 100
\t}
\tlogger := pkglogger.Init(pkglogger.Config{
\t\tLevel:      level,
\t\tFilePath:   cfg.Path,
\t\tMaxSizeMB:  maxSizeMB,
\t\tMaxBackups: cfg.MaxBackups,
\t\tMaxAgeDays: cfg.MaxAgeDays,
\t\tService:    service,
\t})
\tlog.Logger = logger
\tzerolog.DefaultContextLogger = &logger
\tif os.Getenv("ENV") != "production" {
\t\tconsoleWriter := zerolog.ConsoleWriter{Out: os.Stderr}
\t\tlog.Logger = zerolog.New(zerolog.MultiLevelWriter(consoleWriter, logger)).
\t\t\tWith().Timestamp().Str("service", service).Logger()
\t}
\treturn logger
}
'''

metrics_go = '''package internal

import "github.com/prometheus/client_golang/prometheus"

// NewRegistry creates a new Prometheus registry with default process + Go collectors.
func NewRegistry() *prometheus.Registry {
\treg := prometheus.NewRegistry()
\treg.MustRegister(
\t\tprometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
\t\tprometheus.NewGoCollector(),
\t)
\treturn reg
}
'''

files = {
    'database.go': database_go,
    'service.go':  service_go,
    'handlers.go': handlers_go,
    'routes.go':   routes_go,
    'errors.go':   errors_go,
    'logger.go':   logger_go,
    'metrics.go':  metrics_go,
}

for name, content in files.items():
    path = os.path.join(base, name)
    with open(path, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f'wrote {name} ({len(content)} bytes)')

print('all done')
