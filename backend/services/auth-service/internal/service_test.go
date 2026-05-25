package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// testAuthService builds a minimal *AuthService suitable for unit tests that
// do not need a database, Kafka, or real JWT keys.
func testAuthService() *AuthService {
	return &AuthService{
		cfg: &Config{
			OAuth: OAuthCfg{
				Google: OAuthProviderCfg{
					ClientID:     "test-google-client-id",
					ClientSecret: "test-google-client-secret",
					RedirectURL:  "http://localhost:8081/auth/google/callback",
				},
				GitHub: OAuthProviderCfg{
					ClientID:     "test-github-client-id",
					ClientSecret: "test-github-client-secret",
					RedirectURL:  "http://localhost:8081/auth/github/callback",
				},
			},
		},
		googleCfg: &oauth2.Config{
			ClientID:     "test-google-client-id",
			ClientSecret: "test-google-client-secret",
			RedirectURL:  "http://localhost:8081/auth/google/callback",
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		},
		githubCfg: &oauth2.Config{
			ClientID:     "test-github-client-id",
			ClientSecret: "test-github-client-secret",
			RedirectURL:  "http://localhost:8081/auth/github/callback",
			Scopes:       []string{"read:user", "user:email"},
			Endpoint:     github.Endpoint,
		},
	}
}

// TestBuildAuthURL_ValidProviders checks that Google and GitHub both produce
// a non-empty URL that contains the generated state parameter.
func TestBuildAuthURL_ValidProviders(t *testing.T) {
	svc := testAuthService()

	providers := []string{"google", "github"}
	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			authURL, state, verifier, err := svc.BuildAuthURL(provider)
			require.NoError(t, err)
			assert.NotEmpty(t, authURL, "URL must not be empty")
			assert.NotEmpty(t, state, "state must not be empty")
			assert.NotEmpty(t, verifier, "PKCE verifier must not be empty")
			assert.Contains(t, authURL, state, "URL should embed the state parameter")
		})
	}
}

// TestBuildAuthURL_InvalidProvider checks that an unknown provider returns an error.
func TestBuildAuthURL_InvalidProvider(t *testing.T) {
	svc := testAuthService()

	_, _, _, err := svc.BuildAuthURL("twitter")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown provider")
}

// TestBuildAuthURL_UniqueStates verifies that two successive calls generate
// different state and verifier values (entropy is real).
func TestBuildAuthURL_UniqueStates(t *testing.T) {
	svc := testAuthService()

	_, state1, verifier1, err := svc.BuildAuthURL("google")
	require.NoError(t, err)
	_, state2, verifier2, err := svc.BuildAuthURL("google")
	require.NoError(t, err)

	assert.NotEqual(t, state1, state2, "states should be unique across calls")
	assert.NotEqual(t, verifier1, verifier2, "PKCE verifiers should be unique across calls")
}

// TestProviderConfig checks that providerConfig returns the correct config.
func TestProviderConfig_Valid(t *testing.T) {
	svc := testAuthService()

	googleCfg, err := svc.providerConfig("google")
	require.NoError(t, err)
	assert.Equal(t, "test-google-client-id", googleCfg.ClientID)

	githubCfg, err := svc.providerConfig("github")
	require.NoError(t, err)
	assert.Equal(t, "test-github-client-id", githubCfg.ClientID)
}

func TestProviderConfig_Invalid(t *testing.T) {
	svc := testAuthService()
	_, err := svc.providerConfig("unknown")
	require.Error(t, err)
}

// TestHashToken verifies SHA-256 hex output properties.
func TestHashToken(t *testing.T) {
	t.Run("deterministic", func(t *testing.T) {
		assert.Equal(t, hashToken("abc"), hashToken("abc"),
			"same input must produce same hash")
	})
	t.Run("collision resistance", func(t *testing.T) {
		assert.NotEqual(t, hashToken("token-a"), hashToken("token-b"),
			"different inputs must produce different hashes")
	})
	t.Run("output length", func(t *testing.T) {
		h := hashToken("some-opaque-token")
		assert.Len(t, h, 64, "SHA-256 hex digest must be 64 characters")
	})
}
