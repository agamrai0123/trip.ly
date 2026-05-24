package internal

import (
	"context"
	"errors"

	"github.com/agamrai0123/wanderplan/pkg/kafka"
	"github.com/rs/zerolog/log"
)

// CollaborationService handles collaboration business logic.
type CollaborationService struct {
	repo     *CollaboratorRepo
	producer *kafka.Producer
}

// NewCollaborationService constructs the service.
func NewCollaborationService(repo *CollaboratorRepo, producer *kafka.Producer) *CollaborationService {
	return &CollaborationService{repo: repo, producer: producer}
}

// List returns all collaborators for a trip.
func (s *CollaborationService) List(ctx context.Context, tripID string) ([]*Collaborator, error) {
	return s.repo.List(ctx, tripID)
}

// Invite adds a new collaborator to a trip.
func (s *CollaborationService) Invite(ctx context.Context, tripID string, req *InviteCollaboratorRequest) (*Collaborator, error) {
	c, err := s.repo.Invite(ctx, tripID, req.UserID, req.Role)
	if err != nil {
		return nil, err
	}
	s.publishEvent(ctx, kafka.TopicCollabEvents, "collaboration.invited", map[string]string{
		"trip_id": tripID,
		"user_id": req.UserID,
		"role":    req.Role,
	})
	return c, nil
}

// Update changes a collaborator's role.
func (s *CollaborationService) Update(ctx context.Context, tripID, userID string, req *UpdateCollaboratorRequest) (*Collaborator, error) {
	c, err := s.repo.Update(ctx, tripID, userID, req)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	return c, err
}

// Remove removes a collaborator.
func (s *CollaborationService) Remove(ctx context.Context, tripID, userID string) error {
	return s.repo.Remove(ctx, tripID, userID)
}

func (s *CollaborationService) publishEvent(ctx context.Context, topic, eventType string, payload any) {
	if s.producer == nil {
		return
	}
	if err := s.producer.Publish(ctx, topic, eventType, payload); err != nil {
		log.Warn().Err(err).Str("topic", topic).Str("event", eventType).Msg("kafka publish failed")
	}
}
