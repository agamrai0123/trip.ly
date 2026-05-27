package internal

import (
	"context"
	"errors"
	"fmt"

	"github.com/agamrai0123/wanderplan/pkg/kafka"
	"github.com/rs/zerolog/log"
)

// TripService orchestrates trip, day and item business logic.
type TripService struct {
	trips    *TripRepo
	days     *DayRepo
	items    *ItemRepo
	producer *kafka.Producer
}

// NewTripService constructs the service with all dependencies.
func NewTripService(trips *TripRepo, days *DayRepo, items *ItemRepo, producer *kafka.Producer) *TripService {
	return &TripService{trips: trips, days: days, items: items, producer: producer}
}

// ListTrips returns all trips for a user.
func (s *TripService) ListTrips(ctx context.Context, userID string) ([]*Trip, error) {
	return s.trips.List(ctx, userID)
}

// GetTrip fetches a trip and verifies the caller has access.
func (s *TripService) GetTrip(ctx context.Context, id, userID string) (*Trip, error) {
	t, err := s.trips.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t.UserID != userID && t.Visibility != "shared" {
		return nil, ErrNotFound
	}
	return t, nil
}

// CreateTrip creates a new trip for the user.
func (s *TripService) CreateTrip(ctx context.Context, userID string, req *CreateTripRequest) (*Trip, error) {
	t := &Trip{
		UserID:        userID,
		Title:         req.Title,
		Destination:   req.Destination,
		CoverImageURL: req.CoverImageURL,
		Status:        coalesce(req.Status, "draft"),
		Visibility:    coalesce(req.Visibility, "private"),
		BudgetTotal:   req.BudgetTotal,
		Currency:      coalesce(req.Currency, "USD"),
	}
	if req.StartDate != nil {
		if parsed, err := parseDate(*req.StartDate); err == nil {
			t.StartDate = &parsed
		}
	}
	if req.EndDate != nil {
		if parsed, err := parseDate(*req.EndDate); err == nil {
			t.EndDate = &parsed
		}
	}
	created, err := s.trips.Create(ctx, t)
	if err != nil {
		return nil, err
	}
	s.publishEvent(ctx, kafka.TopicTripEvents, "trip.created", map[string]string{
		"trip_id": created.ID,
		"user_id": userID,
	})
	return created, nil
}

// UpdateTrip applies a partial update to a trip.
func (s *TripService) UpdateTrip(ctx context.Context, id, userID string, req *UpdateTripRequest) (*Trip, error) {
	t, err := s.trips.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if t.UserID != userID {
		return nil, ErrNotFound
	}
	updated, err := s.trips.Update(ctx, id, req)
	if err != nil {
		return nil, err
	}
	s.publishEvent(ctx, kafka.TopicTripEvents, "trip.updated", map[string]string{
		"trip_id": id,
		"user_id": userID,
	})
	return updated, nil
}

// DeleteTrip removes a trip owned by the user.
func (s *TripService) DeleteTrip(ctx context.Context, id, userID string) error {
	t, err := s.trips.GetByID(ctx, id)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if t.UserID != userID {
		return ErrNotFound
	}
	if err := s.trips.Delete(ctx, id); err != nil {
		return err
	}
	s.publishEvent(ctx, kafka.TopicTripEvents, "trip.deleted", map[string]string{
		"trip_id": id,
		"user_id": userID,
	})
	return nil
}

// DuplicateTrip creates a copy of an existing trip (without days/items for now).
func (s *TripService) DuplicateTrip(ctx context.Context, id, userID string) (*Trip, error) {
	src, err := s.GetTrip(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	req := &CreateTripRequest{
		Title:         fmt.Sprintf("Copy of %s", src.Title),
		Destination:   src.Destination,
		CoverImageURL: src.CoverImageURL,
		Status:        "draft",
		Visibility:    "private",
		BudgetTotal:   src.BudgetTotal,
		Currency:      src.Currency,
	}
	return s.CreateTrip(ctx, userID, req)
}

// GetStats returns aggregate statistics for the dashboard.
func (s *TripService) GetStats(ctx context.Context, userID string) (*TripStats, error) {
	return s.trips.Stats(ctx, userID)
}

// ── Day operations ──────────────────────────────────────────

// CreateDay adds a new itinerary day to a trip.
func (s *TripService) CreateDay(ctx context.Context, tripID, userID string, req *CreateDayRequest) (*ItineraryDay, error) {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return nil, err
	}
	d := &ItineraryDay{
		TripID:    tripID,
		DayNumber: req.DayNumber,
		Notes:     req.Notes,
	}
	if req.Date != nil {
		if parsed, err := parseDate(*req.Date); err == nil {
			d.Date = &parsed
		}
	}
	return s.days.Create(ctx, d)
}

// ListDays returns the days for a trip.
func (s *TripService) ListDays(ctx context.Context, tripID, userID string) ([]*ItineraryDay, error) {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return nil, err
	}
	return s.days.ListByTrip(ctx, tripID)
}

// UpdateDay patches a day.
func (s *TripService) UpdateDay(ctx context.Context, tripID, dayID, userID string, req *UpdateDayRequest) (*ItineraryDay, error) {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return nil, err
	}
	return s.days.Update(ctx, dayID, req)
}

// DeleteDay removes a day.
func (s *TripService) DeleteDay(ctx context.Context, tripID, dayID, userID string) error {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return err
	}
	return s.days.Delete(ctx, dayID)
}

// ── Item operations ─────────────────────────────────────────

// CreateItem adds a new item to a day.
func (s *TripService) CreateItem(ctx context.Context, tripID, dayID, userID string, req *CreateItemRequest) (*ItineraryItem, error) {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return nil, err
	}
	it := &ItineraryItem{
		DayID:       dayID,
		TripID:      tripID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		PlaceID:     req.PlaceID,
		Type:        req.Type,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Cost:        req.Cost,
		Currency:    coalesce(req.Currency, "USD"),
		OrderIndex:  req.OrderIndex,
	}
	return s.items.Create(ctx, it)
}

// ListItems returns items for a day.
func (s *TripService) ListItems(ctx context.Context, tripID, dayID, userID string) ([]*ItineraryItem, error) {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return nil, err
	}
	return s.items.ListByDay(ctx, dayID)
}

// UpdateItem patches an item.
func (s *TripService) UpdateItem(ctx context.Context, tripID, itemID, userID string, req *UpdateItemRequest) (*ItineraryItem, error) {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return nil, err
	}
	return s.items.Update(ctx, itemID, req)
}

// DeleteItem removes an item.
func (s *TripService) DeleteItem(ctx context.Context, tripID, itemID, userID string) error {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return err
	}
	return s.items.Delete(ctx, itemID)
}

// ReorderItems bulk-updates order_index values for dnd-kit reorder events.
func (s *TripService) ReorderItems(ctx context.Context, tripID, userID string, req *ReorderRequest) error {
	if _, err := s.GetTrip(ctx, tripID, userID); err != nil {
		return err
	}
	return s.items.Reorder(ctx, tripID, req.Items)
}

// ── Helpers ──────────────────────────────────────────────────

func (s *TripService) publishEvent(ctx context.Context, topic, eventType string, payload any) {
	if s.producer == nil {
		return
	}
	if err := s.producer.Publish(ctx, topic, eventType, payload); err != nil {
		log.Warn().Err(err).Str("topic", topic).Str("event", eventType).Msg("kafka publish failed")
	}
}
