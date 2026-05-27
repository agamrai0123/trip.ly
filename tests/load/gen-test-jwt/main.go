// gen-test-jwt generates a signed RS256 JWT for load-testing authenticated endpoints.
// Usage: go run ./tests/load/gen-test-jwt/
// Reads JWT_PRIVATE_KEY (base64 PEM) from .env or environment.
// Prints a valid Bearer token to stdout.
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	privB64 := os.Getenv("JWT_PRIVATE_KEY")
	if privB64 == "" {
		// Try reading from .env in repo root
		data, err := os.ReadFile(".env")
		if err != nil {
			fmt.Fprintln(os.Stderr, "JWT_PRIVATE_KEY not set and .env not found")
			os.Exit(1)
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimRight(line, "\r")
			if strings.HasPrefix(line, "JWT_PRIVATE_KEY=") {
				privB64 = strings.TrimPrefix(line, "JWT_PRIVATE_KEY=")
				break
			}
		}
	}
	if privB64 == "" {
		fmt.Fprintln(os.Stderr, "JWT_PRIVATE_KEY not found in environment or .env")
		os.Exit(1)
	}

	pemBytes, err := base64.StdEncoding.DecodeString(privB64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "base64 decode:", err)
		os.Exit(1)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(pemBytes)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse RSA key:", err)
		os.Exit(1)
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":    "00000000-0000-0000-0000-000000000001",
		"email":      "loadtest@wanderplan.local",
		"name":       "Load Test User",
		"avatar_url": "",
		"iat":        now.Unix(),
		"exp":        now.Add(24 * time.Hour).Unix(), // long expiry for load tests
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signed, err := token.SignedString(key)
	if err != nil {
		fmt.Fprintln(os.Stderr, "sign token:", err)
		os.Exit(1)
	}

	fmt.Print(signed)
}
