package internal

import (
	"context"
	"errors"

	proto "github.com/agamrai0123/wanderplan/proto/gen/wanderplan/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// UserService handles user business logic.
type UserService struct {
	repo        *UserRepo
	tripSvcAddr string
}

// NewUserService constructs the service.
func NewUserService(repo *UserRepo, tripSvcAddr string) *UserService {
	return &UserService{repo: repo, tripSvcAddr: tripSvcAddr}
}

// GetMe returns the current user's profile.
func (s *UserService) GetMe(ctx context.Context, userID string) (*User, error) {
	return s.repo.GetByID(ctx, userID)
}

// UpdateMe applies a partial update to the current user's profile.
func (s *UserService) UpdateMe(ctx context.Context, userID string, req *UpdateUserRequest) (*User, error) {
	u, err := s.repo.Update(ctx, userID, req)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	return u, err
}

// GetMyTrips delegates to trip-service via gRPC.
func (s *UserService) GetMyTrips(ctx context.Context, userID string) (*proto.ListTripsByUserResponse, error) {
	conn, err := grpc.NewClient(s.tripSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close() //nolint:errcheck

	client := proto.NewTripServiceClient(conn)
	return client.ListTripsByUser(ctx, &proto.ListTripsByUserRequest{UserId: userID})
}

// GetMyStats delegates stats computation to trip-service via gRPC.
func (s *UserService) GetMyStats(ctx context.Context, userID string) (*proto.GetTripStatsResponse, error) {
	conn, err := grpc.NewClient(s.tripSvcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close() //nolint:errcheck

	client := proto.NewTripServiceClient(conn)
	return client.GetTripStats(ctx, &proto.GetTripStatsRequest{UserId: userID})
}
