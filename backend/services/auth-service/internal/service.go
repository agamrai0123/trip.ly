package internal

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"github.com/agamrai0123/wanderplan/pkg/jwt"
	"github.com/agamrai0123/wanderplan/pkg/kafka"
)

const (
	refreshTokenTTL = 30 * 24 * time.Hour
	accessTokenTTL  = 15 * time.Minute
)

// AuthService orchestrates OAuth PKCE, JWT issuance, and Kafka events.
type AuthService struct {
	cfg           *Config
	users         *UserRepo
	refreshTokens *RefreshTokenRepo
	jwtMgr        *jwt.Manager
	producer      *kafka.Producer
	googleCfg     *oauth2.Config
	githubCfg     *oauth2.Config
}

// NewAuthService constructs the service with all dependencies.
func NewAuthService(
	cfg *Config,
	users *UserRepo,
	rt *RefreshTokenRepo,
	jwtMgr *jwt.Manager,
	producer *kafka.Producer,
) *AuthService {
	gCfg := &oauth2.Config{
		ClientID:     cfg.OAuth.Google.ClientID,
		ClientSecret: cfg.OAuth.Google.ClientSecret,
		RedirectURL:  cfg.OAuth.Google.RedirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}
	ghCfg := &oauth2.Config{
		ClientID:     cfg.OAuth.GitHub.ClientID,
		ClientSecret: cfg.OAuth.GitHub.ClientSecret,
		RedirectURL:  cfg.OAuth.GitHub.RedirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}
	return &AuthService{
		cfg:           cfg,
		users:         users,
		refreshTokens: rt,
		jwtMgr:        jwtMgr,
		producer:      producer,
		googleCfg:     gCfg,
		githubCfg:     ghCfg,
	}
}

// BuildAuthURL generates the OAuth2 PKCE authorization URL.
// Returns (url, state, codeVerifier, error).
func (s *AuthService) BuildAuthURL(provider string) (string, string, string, error) {
	state, err := randomBase64(32)
	if err != nil {
		return "", "", "", fmt.Errorf("generate state: %w", err)
	}
	verifier, err := randomBase64(32)
	if err != nil {
		return "", "", "", fmt.Errorf("generate verifier: %w", err)
	}
	cfg, err := s.providerConfig(provider)
	if err != nil {
		return "", "", "", err
	}
	url := cfg.AuthCodeURL(state,
		oauth2.AccessTypeOffline,
		oauth2.S256ChallengeOption(verifier),
	)
	return url, state, verifier, nil
}

// ExchangeCode swaps the authorization code for user info and issues WanderPlan tokens.
func (s *AuthService) ExchangeCode(
	ctx context.Context,
	provider, code, verifier string,
) (*LoginResponse, string, error) {
	cfg, err := s.providerConfig(provider)
	if err != nil {
		return nil, "", err
	}
	token, err := cfg.Exchange(ctx, code, oauth2.VerifierOption(verifier))
	if err != nil {
		return nil, "", fmt.Errorf("oauth exchange: %w", err)
	}
	userInfo, err := s.fetchUserInfo(ctx, provider, token)
	if err != nil {
		return nil, "", err
	}
	user, err := s.users.Upsert(ctx, userInfo)
	if err != nil {
		return nil, "", fmt.Errorf("upsert user: %w", err)
	}
	accessToken, err := s.jwtMgr.Sign(user.ID, user.Email, user.Name, user.AvatarURL)
	if err != nil {
		return nil, "", fmt.Errorf("sign jwt: %w", err)
	}
	rawRefresh, err := s.refreshTokens.Create(ctx, user.ID, refreshTokenTTL)
	if err != nil {
		return nil, "", fmt.Errorf("create refresh token: %w", err)
	}
	payload, _ := json.Marshal(map[string]string{
		"user_id": user.ID,
		"email":   user.Email,
		"event":   "user.login",
	})
	if err := s.producer.Publish(ctx, kafka.TopicAuthEvents, "user.login", payload); err != nil {
		log.Warn().Err(err).Msg("failed to publish auth event")
	}
	return &LoginResponse{AccessToken: accessToken, User: user}, rawRefresh, nil
}

// Refresh validates a refresh token, rotates it, and issues a new access token.
func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*RefreshResponse, string, error) {
	userID, newRaw, err := s.refreshTokens.Rotate(ctx, rawToken, refreshTokenTTL)
	if err != nil {
		return nil, "", err
	}
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, "", err
	}
	accessToken, err := s.jwtMgr.Sign(user.ID, user.Email, user.Name, user.AvatarURL)
	if err != nil {
		return nil, "", err
	}
	return &RefreshResponse{AccessToken: accessToken}, newRaw, nil
}

// Logout revokes all refresh tokens for the user.
func (s *AuthService) Logout(ctx context.Context, userID string) error {
	return s.refreshTokens.DeleteByUserID(ctx, userID)
}

// ValidateToken validates a JWT and returns the claims.
func (s *AuthService) ValidateToken(tokenStr string) (*jwt.Claims, error) {
	return s.jwtMgr.Parse(tokenStr)
}

func (s *AuthService) providerConfig(provider string) (*oauth2.Config, error) {
	switch provider {
	case "google":
		return s.googleCfg, nil
	case "github":
		return s.githubCfg, nil
	default:
		return nil, fmt.Errorf("unknown provider: %s", provider)
	}
}

type googleUserInfo struct {
	Sub     string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type githubUserInfo struct {
	Login     string `json:"login"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

func (s *AuthService) fetchUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*User, error) {
	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	switch provider {
	case "google":
		resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("google userinfo status %d", resp.StatusCode)
		}
		var info googleUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return nil, err
		}
		return &User{ID: info.Sub, Email: info.Email, Name: info.Name, AvatarURL: info.Picture, Provider: "google"}, nil
	case "github":
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("github user status %d", resp.StatusCode)
		}
		var info githubUserInfo
		if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return nil, err
		}
		if info.Email == "" {
			info.Email, _ = fetchGitHubPrimaryEmail(client)
		}
		name := info.Name
		if name == "" {
			name = info.Login
		}
		return &User{ID: fmt.Sprintf("gh_%s", info.Login), Email: info.Email, Name: name, AvatarURL: info.AvatarURL, Provider: "github"}, nil
	}
	return nil, fmt.Errorf("unknown provider: %s", provider)
}

type ghEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func fetchGitHubPrimaryEmail(client *http.Client) (string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var emails []ghEmail
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", err
	}
	for _, e := range emails {
		if e.Primary {
			return e.Email, nil
		}
	}
	return "", nil
}

func randomBase64(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
