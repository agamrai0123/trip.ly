package internal

import (
	"context"
	"errors"
)

// NotificationService handles notification business logic.
type NotificationService struct {
	repo *NotificationRepo
}

// NewNotificationService constructs the service.
func NewNotificationService(repo *NotificationRepo) *NotificationService {
	return &NotificationService{repo: repo}
}

// GetByUser returns all notifications for a user.
func (s *NotificationService) GetByUser(ctx context.Context, userID string) ([]*Notification, error) {
	return s.repo.GetByUser(ctx, userID)
}

// MarkRead marks a single notification as read.
func (s *NotificationService) MarkRead(ctx context.Context, id, userID string) error {
	err := s.repo.MarkRead(ctx, id, userID)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}
	return err
}

// MarkAllRead marks all notifications for a user as read.
func (s *NotificationService) MarkAllRead(ctx context.Context, userID string) error {
	return s.repo.MarkAllRead(ctx, userID)
}
