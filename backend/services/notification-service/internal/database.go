package internal

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a row does not exist.
var ErrNotFound = errors.New("not found")

// NotificationRepo handles notification persistence.
type NotificationRepo struct{ pool *pgxpool.Pool }

// NewNotificationRepo creates a repo backed by the given pool.
func NewNotificationRepo(pool *pgxpool.Pool) *NotificationRepo {
	return &NotificationRepo{pool: pool}
}

// GetByUser returns all notifications for a user, newest first.
func (r *NotificationRepo) GetByUser(ctx context.Context, userID string) ([]*Notification, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, type, title, body, read_at, created_at
		FROM wanderplan.notifications WHERE user_id=$1 ORDER BY created_at DESC LIMIT 100`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ns []*Notification
	for rows.Next() {
		n := &Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.ReadAt, &n.CreatedAt); err != nil {
			return nil, err
		}
		ns = append(ns, n)
	}
	return ns, rows.Err()
}

// Create persists a new notification and returns it.
func (r *NotificationRepo) Create(ctx context.Context, n *Notification) (*Notification, error) {
	n.ID = uuid.New().String()
	n.CreatedAt = time.Now().UTC()
	row := r.pool.QueryRow(ctx, `
		INSERT INTO wanderplan.notifications (id, user_id, type, title, body, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, user_id, type, title, body, read_at, created_at`,
		n.ID, n.UserID, n.Type, n.Title, n.Body, n.CreatedAt)
	out := &Notification{}
	err := row.Scan(&out.ID, &out.UserID, &out.Type, &out.Title, &out.Body, &out.ReadAt, &out.CreatedAt)
	return out, err
}

// MarkRead sets read_at on a single notification owned by userID.
func (r *NotificationRepo) MarkRead(ctx context.Context, id, userID string) error {
	tag, err := r.pool.Exec(ctx, `
		UPDATE wanderplan.notifications SET read_at=NOW()
		WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// MarkAllRead marks every unread notification for a user as read.
func (r *NotificationRepo) MarkAllRead(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE wanderplan.notifications SET read_at=NOW()
		WHERE user_id=$1 AND read_at IS NULL`, userID)
	return err
}

// GetByID fetches a single notification, checking pgx.ErrNoRows.
func (r *NotificationRepo) GetByID(ctx context.Context, id string) (*Notification, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, type, title, body, read_at, created_at
		FROM wanderplan.notifications WHERE id=$1`, id)
	n := &Notification{}
	err := row.Scan(&n.ID, &n.UserID, &n.Type, &n.Title, &n.Body, &n.ReadAt, &n.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return n, err
}
