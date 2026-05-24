// Package grpc provides a gRPC client factory with connection pooling and retry.
package grpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// ConnConfig holds options for dialling a gRPC endpoint.
type ConnConfig struct {
	// Target is "host:port".
	Target string
	// Timeout is the dial timeout. Defaults to 5s.
	Timeout time.Duration
}

// Dial creates a gRPC client connection to the given target.
// It uses insecure credentials (mTLS is handled at the infrastructure layer).
// The connection uses wait-for-ready semantics and exponential backoff.
func Dial(ctx context.Context, cfg ConnConfig) (*grpc.ClientConn, error) {
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	dialCtx, cancel := context.WithTimeout(ctx, cfg.Timeout)
	defer cancel()

	conn, err := grpc.DialContext( //nolint:staticcheck // DialContext deprecated in 1.63 but still valid
		dialCtx,
		cfg.Target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  100 * time.Millisecond,
				Multiplier: 1.6,
				Jitter:     0.2,
				MaxDelay:   5 * time.Second,
			},
			MinConnectTimeout: 2 * time.Second,
		}),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                30 * time.Second,
			Timeout:             10 * time.Second,
			PermitWithoutStream: true,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("grpc dial %s: %w", cfg.Target, err)
	}
	return conn, nil
}
