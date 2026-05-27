package internal

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a queried row does not exist.
var ErrNotFound = errors.New("not found")

// TripRepo handles trip persistence against PostgreSQL.
type TripRepo struct{ pool *pgxpool.Pool }

// NewTripRepo creates a TripRepo backed by the supplied pool.
func NewTripRepo(pool *pgxpool.Pool) *TripRepo { return &TripRepo{pool: pool} }

// List returns all trips for the given user ordered by created_at desc.
func (r *TripRepo) List(ctx context.Context, userID string) ([]*Trip, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, title, destination, cover_image_url,
		       start_date, end_date, status, visibility,
		       budget_total, currency, created_at, updated_at
		FROM wanderplan.trips
		WHERE user_id = $1
		ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*Trip
	for rows.Next() {
		t := &Trip{}
		if err := rows.Scan(
			&t.ID, &t.UserID, &t.Title, &t.Destination, &t.CoverImageURL,
			&t.StartDate, &t.EndDate, &t.Status, &t.Visibility,
			&t.BudgetTotal, &t.Currency, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		trips = append(trips, t)
	}
	return trips, rows.Err()
}

// GetByID fetches a single trip by ID.
func (r *TripRepo) GetByID(ctx context.Context, id string) (*Trip, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, user_id, title, destination, cover_image_url,
		       start_date, end_date, status, visibility,
		       budget_total, currency, created_at, updated_at
		FROM wanderplan.trips WHERE id = $1`, id)
	t := &Trip{}
	err := row.Scan(
		&t.ID, &t.UserID, &t.Title, &t.Destination, &t.CoverImageURL,
		&t.StartDate, &t.EndDate, &t.Status, &t.Visibility,
		&t.BudgetTotal, &t.Currency, &t.CreatedAt, &t.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return t, err
}

// Create inserts a new trip and returns it.
func (r *TripRepo) Create(ctx context.Context, t *Trip) (*Trip, error) {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now().UTC()
	t.UpdatedAt = t.CreatedAt
	row := r.pool.QueryRow(ctx, `
		INSERT INTO wanderplan.trips
		  (id, user_id, title, destination, cover_image_url,
		   start_date, end_date, status, visibility, budget_total, currency, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING id, user_id, title, destination, cover_image_url,
		          start_date, end_date, status, visibility, budget_total, currency, created_at, updated_at`,
		t.ID, t.UserID, t.Title, t.Destination, t.CoverImageURL,
		t.StartDate, t.EndDate, t.Status, t.Visibility, t.BudgetTotal, t.Currency,
		t.CreatedAt, t.UpdatedAt,
	)
	out := &Trip{}
	err := row.Scan(
		&out.ID, &out.UserID, &out.Title, &out.Destination, &out.CoverImageURL,
		&out.StartDate, &out.EndDate, &out.Status, &out.Visibility,
		&out.BudgetTotal, &out.Currency, &out.CreatedAt, &out.UpdatedAt,
	)
	return out, err
}

// Update applies partial changes to an existing trip.
func (r *TripRepo) Update(ctx context.Context, id string, req *UpdateTripRequest) (*Trip, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Destination != nil {
		existing.Destination = *req.Destination
	}
	if req.CoverImageURL != nil {
		existing.CoverImageURL = *req.CoverImageURL
	}
	if req.Status != nil {
		existing.Status = *req.Status
	}
	if req.Visibility != nil {
		existing.Visibility = *req.Visibility
	}
	if req.BudgetTotal != nil {
		existing.BudgetTotal = *req.BudgetTotal
	}
	if req.Currency != nil {
		existing.Currency = *req.Currency
	}
	existing.UpdatedAt = time.Now().UTC()

	row := r.pool.QueryRow(ctx, `
		UPDATE wanderplan.trips
		SET title=$2, destination=$3, cover_image_url=$4,
		    status=$5, visibility=$6, budget_total=$7, currency=$8, updated_at=$9
		WHERE id=$1
		RETURNING id, user_id, title, destination, cover_image_url,
		          start_date, end_date, status, visibility, budget_total, currency, created_at, updated_at`,
		existing.ID, existing.Title, existing.Destination, existing.CoverImageURL,
		existing.Status, existing.Visibility, existing.BudgetTotal, existing.Currency, existing.UpdatedAt,
	)
	out := &Trip{}
	err = row.Scan(
		&out.ID, &out.UserID, &out.Title, &out.Destination, &out.CoverImageURL,
		&out.StartDate, &out.EndDate, &out.Status, &out.Visibility,
		&out.BudgetTotal, &out.Currency, &out.CreatedAt, &out.UpdatedAt,
	)
	return out, err
}

// Delete removes a trip and all its days/items via CASCADE.
func (r *TripRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM wanderplan.trips WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Stats returns aggregate statistics for a user's trips.
func (r *TripRepo) Stats(ctx context.Context, userID string) (*TripStats, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
		  COUNT(*),
		  COUNT(DISTINCT destination),
		  COALESCE(SUM(CASE WHEN start_date IS NOT NULL AND end_date IS NOT NULL
		                    THEN (end_date::date - start_date::date + 1) ELSE 0 END), 0),
		  COALESCE(SUM(budget_total), 0)
		FROM wanderplan.trips WHERE user_id=$1`, userID)
	s := &TripStats{}
	return s, row.Scan(&s.TotalTrips, &s.TotalCountries, &s.TotalDays, &s.TotalBudget)
}

// ─────────────────────────────────────────────────────────────
// DayRepo
// ─────────────────────────────────────────────────────────────

// DayRepo handles itinerary_days persistence.
type DayRepo struct{ pool *pgxpool.Pool }

// NewDayRepo creates a DayRepo backed by the supplied pool.
func NewDayRepo(pool *pgxpool.Pool) *DayRepo { return &DayRepo{pool: pool} }

// ListByTrip returns all days for the given trip ordered by day_number.
func (r *DayRepo) ListByTrip(ctx context.Context, tripID string) ([]*ItineraryDay, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, trip_id, day_number, date, notes
		FROM wanderplan.itinerary_days
		WHERE trip_id=$1 ORDER BY day_number`, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var days []*ItineraryDay
	for rows.Next() {
		d := &ItineraryDay{}
		if err := rows.Scan(&d.ID, &d.TripID, &d.DayNumber, &d.Date, &d.Notes); err != nil {
			return nil, err
		}
		days = append(days, d)
	}
	return days, rows.Err()
}

// Create inserts a new itinerary day.
func (r *DayRepo) Create(ctx context.Context, d *ItineraryDay) (*ItineraryDay, error) {
	d.ID = uuid.New().String()
	row := r.pool.QueryRow(ctx, `
		INSERT INTO wanderplan.itinerary_days (id, trip_id, day_number, date, notes)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, trip_id, day_number, date, notes`,
		d.ID, d.TripID, d.DayNumber, d.Date, d.Notes)
	out := &ItineraryDay{}
	err := row.Scan(&out.ID, &out.TripID, &out.DayNumber, &out.Date, &out.Notes)
	return out, err
}

// Update applies partial changes to a day.
func (r *DayRepo) Update(ctx context.Context, id string, req *UpdateDayRequest) (*ItineraryDay, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Notes != nil {
		existing.Notes = *req.Notes
	}
	row := r.pool.QueryRow(ctx, `
		UPDATE wanderplan.itinerary_days SET notes=$2 WHERE id=$1
		RETURNING id, trip_id, day_number, date, notes`, existing.ID, existing.Notes)
	out := &ItineraryDay{}
	err = row.Scan(&out.ID, &out.TripID, &out.DayNumber, &out.Date, &out.Notes)
	return out, err
}

// GetByID fetches a single day.
func (r *DayRepo) GetByID(ctx context.Context, id string) (*ItineraryDay, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, trip_id, day_number, date, notes
		FROM wanderplan.itinerary_days WHERE id=$1`, id)
	d := &ItineraryDay{}
	err := row.Scan(&d.ID, &d.TripID, &d.DayNumber, &d.Date, &d.Notes)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return d, err
}

// Delete removes a day and cascades to its items.
func (r *DayRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM wanderplan.itinerary_days WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// ─────────────────────────────────────────────────────────────
// ItemRepo
// ─────────────────────────────────────────────────────────────

// ItemRepo handles itinerary_items persistence.
type ItemRepo struct{ pool *pgxpool.Pool }

// NewItemRepo creates an ItemRepo backed by the supplied pool.
func NewItemRepo(pool *pgxpool.Pool) *ItemRepo { return &ItemRepo{pool: pool} }

// ListByDay returns all items for a day ordered by order_index.
func (r *ItemRepo) ListByDay(ctx context.Context, dayID string) ([]*ItineraryItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, day_id, trip_id, title, description, location, place_id,
		       type, start_time, end_time, cost, currency, order_index, created_at, updated_at
		FROM wanderplan.itinerary_items
		WHERE day_id=$1 ORDER BY order_index`, dayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ItineraryItem
	for rows.Next() {
		it := &ItineraryItem{}
		if err := rows.Scan(
			&it.ID, &it.DayID, &it.TripID, &it.Title, &it.Description,
			&it.Location, &it.PlaceID, &it.Type, &it.StartTime, &it.EndTime,
			&it.Cost, &it.Currency, &it.OrderIndex, &it.CreatedAt, &it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, it)
	}
	return items, rows.Err()
}

// Create inserts a new itinerary item.
func (r *ItemRepo) Create(ctx context.Context, it *ItineraryItem) (*ItineraryItem, error) {
	it.ID = uuid.New().String()
	it.CreatedAt = time.Now().UTC()
	it.UpdatedAt = it.CreatedAt
	row := r.pool.QueryRow(ctx, `
		INSERT INTO wanderplan.itinerary_items
		  (id, day_id, trip_id, title, description, location, place_id,
		   type, start_time, end_time, cost, currency, order_index, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)
		RETURNING id, day_id, trip_id, title, description, location, place_id,
		          type, start_time, end_time, cost, currency, order_index, created_at, updated_at`,
		it.ID, it.DayID, it.TripID, it.Title, it.Description, it.Location, it.PlaceID,
		it.Type, it.StartTime, it.EndTime, it.Cost, it.Currency, it.OrderIndex,
		it.CreatedAt, it.UpdatedAt,
	)
	out := &ItineraryItem{}
	err := row.Scan(
		&out.ID, &out.DayID, &out.TripID, &out.Title, &out.Description,
		&out.Location, &out.PlaceID, &out.Type, &out.StartTime, &out.EndTime,
		&out.Cost, &out.Currency, &out.OrderIndex, &out.CreatedAt, &out.UpdatedAt,
	)
	return out, err
}

// GetByID fetches a single item.
func (r *ItemRepo) GetByID(ctx context.Context, id string) (*ItineraryItem, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, day_id, trip_id, title, description, location, place_id,
		       type, start_time, end_time, cost, currency, order_index, created_at, updated_at
		FROM wanderplan.itinerary_items WHERE id=$1`, id)
	it := &ItineraryItem{}
	err := row.Scan(
		&it.ID, &it.DayID, &it.TripID, &it.Title, &it.Description,
		&it.Location, &it.PlaceID, &it.Type, &it.StartTime, &it.EndTime,
		&it.Cost, &it.Currency, &it.OrderIndex, &it.CreatedAt, &it.UpdatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return it, err
}

// Update applies partial changes to an item.
func (r *ItemRepo) Update(ctx context.Context, id string, req *UpdateItemRequest) (*ItineraryItem, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Title != nil {
		existing.Title = *req.Title
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Location != nil {
		existing.Location = *req.Location
	}
	if req.PlaceID != nil {
		existing.PlaceID = *req.PlaceID
	}
	if req.Type != nil {
		existing.Type = *req.Type
	}
	if req.StartTime != nil {
		existing.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		existing.EndTime = *req.EndTime
	}
	if req.Cost != nil {
		existing.Cost = *req.Cost
	}
	if req.Currency != nil {
		existing.Currency = *req.Currency
	}
	if req.OrderIndex != nil {
		existing.OrderIndex = *req.OrderIndex
	}
	existing.UpdatedAt = time.Now().UTC()

	row := r.pool.QueryRow(ctx, `
		UPDATE wanderplan.itinerary_items
		SET title=$2, description=$3, location=$4, place_id=$5, type=$6,
		    start_time=$7, end_time=$8, cost=$9, currency=$10, order_index=$11, updated_at=$12
		WHERE id=$1
		RETURNING id, day_id, trip_id, title, description, location, place_id,
		          type, start_time, end_time, cost, currency, order_index, created_at, updated_at`,
		existing.ID, existing.Title, existing.Description, existing.Location, existing.PlaceID,
		existing.Type, existing.StartTime, existing.EndTime, existing.Cost, existing.Currency,
		existing.OrderIndex, existing.UpdatedAt,
	)
	out := &ItineraryItem{}
	err = row.Scan(
		&out.ID, &out.DayID, &out.TripID, &out.Title, &out.Description,
		&out.Location, &out.PlaceID, &out.Type, &out.StartTime, &out.EndTime,
		&out.Cost, &out.Currency, &out.OrderIndex, &out.CreatedAt, &out.UpdatedAt,
	)
	return out, err
}

// Delete removes an item by ID.
func (r *ItemRepo) Delete(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM wanderplan.itinerary_items WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// Reorder bulk-updates order_index values for a set of items within a trip.
func (r *ItemRepo) Reorder(ctx context.Context, tripID string, items []ReorderItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx) //nolint:errcheck

	for _, item := range items {
		_, err := tx.Exec(ctx, `
			UPDATE wanderplan.itinerary_items
			SET order_index=$1, updated_at=NOW()
			WHERE id=$2 AND trip_id=$3`,
			item.OrderIndex, item.ID, tripID)
		if err != nil {
			return err
		}
	}
	return tx.Commit(ctx)
}
