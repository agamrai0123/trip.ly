package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/agamrai0123/wanderplan/pkg/errors"
	"github.com/agamrai0123/wanderplan/pkg/middleware"
	"github.com/agamrai0123/wanderplan/pkg/response"
)

// Handlers wires HTTP requests to AuthService methods.
type Handlers struct {
	svc *AuthService
	cfg *Config
}

// NewHandlers constructs the handler set.
func NewHandlers(svc *AuthService, cfg *Config) *Handlers {
	return &Handlers{svc: svc, cfg: cfg}
}

// Login redirects the user to the OAuth provider.
// GET /auth/:provider/login
func (h *Handlers) Login(c *gin.Context) {
	provider := c.Param("provider")
	authURL, state, verifier, err := h.svc.BuildAuthURL(provider)
	if err != nil {
		response.Err(c, errors.BadRequest(err.Error()))
		return
	}
	maxAge := h.cfg.Cookie.StateMaxAgeSecs
	if maxAge == 0 {
		maxAge = 300
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("oauth_state", state, maxAge, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
	c.SetCookie("pkce_verifier", verifier, maxAge, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback handles the OAuth provider redirect.
// GET /auth/:provider/callback
func (h *Handlers) Callback(c *gin.Context) {
	provider := c.Param("provider")
	cookieState, err := c.Cookie("oauth_state")
	if err != nil || cookieState != c.Query("state") {
		response.Err(c, errors.BadRequest("invalid oauth state"))
		return
	}
	verifier, err := c.Cookie("pkce_verifier")
	if err != nil {
		response.Err(c, errors.BadRequest("missing pkce verifier"))
		return
	}
	c.SetCookie("oauth_state", "", -1, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)
	c.SetCookie("pkce_verifier", "", -1, "/", h.cfg.Cookie.Domain, h.cfg.Cookie.Secure, true)

	loginResp, rawRefresh, err := h.svc.ExchangeCode(c.Request.Context(), provider, c.Query("code"), verifier)
	if err != nil {
		response.Err(c, errors.Internal(err.Error()))
		return
	}
	h.setRefreshCookie(c, rawRefresh)
	userJSON, _ := json.Marshal(loginResp.User)
	userB64 := base64.URLEncoding.EncodeToString(userJSON)
	frontendURL := h.cfg.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:5173"
	}
	redirectURL := fmt.Sprintf("%s/auth/callback?access_token=%s&user=%s",
		frontendURL,
		url.QueryEscape(loginResp.AccessToken),
		url.QueryEscape(userB64),
	)
	c.Redirect(http.StatusFound, redirectURL)
}

// Refresh issues a new access token using the refresh_token cookie.
// POST /auth/refresh
func (h *Handlers) Refresh(c *gin.Context) {
	rawToken, err := c.Cookie("refresh_token")
	if err != nil {
		response.Err(c, errors.Unauthorized("missing refresh token"))
		return
	}
	refreshResp, newRaw, err := h.svc.Refresh(c.Request.Context(), rawToken)
	if err != nil {
		h.clearRefreshCookie(c)
		response.Err(c, errors.Unauthorized("invalid or expired refresh token"))
		return
	}
	h.setRefreshCookie(c, newRaw)
	response.OK(c, refreshResp)
}

// Logout revokes the refresh token and clears the cookie.
// POST /auth/logout
func (h *Handlers) Logout(c *gin.Context) {
	userID := middleware.UserID(c)
	if userID != "" {
		_ = h.svc.Logout(c.Request.Context(), userID)
	}
	h.clearRefreshCookie(c)
	response.NoContent(c)
}

// Me returns the authenticated user's claims.
// GET /auth/me
func (h *Handlers) Me(c *gin.Context) {
	userID := middleware.UserID(c)
	email, _ := c.Get("email")
	name, _ := c.Get("name")
	avatar, _ := c.Get("avatar_url")
	response.OK(c, gin.H{
		"user_id":    userID,
		"email":      email,
		"name":       name,
		"avatar_url": avatar,
	})
}

func (h *Handlers) setRefreshCookie(c *gin.Context, raw string) {
	const maxAge = int(30 * 24 * time.Hour / time.Second)
	c.SetSameSite(http.SameSiteNoneMode)
	c.SetCookie("refresh_token", raw, maxAge, "/", h.cfg.Cookie.Domain, true, true)
}

func (h *Handlers) clearRefreshCookie(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/", h.cfg.Cookie.Domain, true, true)
}
