package internal

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	pkgjwt "github.com/agamrai0123/wanderplan/pkg/jwt"
	pbauth "github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"time"
)

// AuthValidator calls auth-service over gRPC to validate a JWT.
type AuthValidator struct {
	client pbauth.AuthServiceClient
	conn   *grpc.ClientConn
}

// NewAuthValidator dials the auth-service gRPC server.
func NewAuthValidator(authAddr string) (*AuthValidator, error) {
	conn, err := grpc.NewClient(authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    30 * time.Second,
			Timeout: 10 * time.Second,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("dial auth-service: %w", err)
	}
	log.Info().Str("addr", authAddr).Msg("connected to auth-service gRPC")
	return &AuthValidator{client: pbauth.NewAuthServiceClient(conn), conn: conn}, nil
}

// Validate calls AuthService.ValidateToken and returns the parsed claims.
func (v *AuthValidator) Validate(ctx context.Context, token string) (*pkgjwt.Claims, error) {
	resp, err := v.client.ValidateToken(ctx, &pbauth.ValidateTokenRequest{Token: token})
	if err != nil {
		return nil, fmt.Errorf("validate token: %w", err)
	}
	return &pkgjwt.Claims{
		UserID:    resp.UserId,
		Email:     resp.Email,
		Name:      resp.Name,
		AvatarURL: resp.AvatarUrl,
	}, nil
}

// Close releases the gRPC connection.
func (v *AuthValidator) Close() error { return v.conn.Close() }
