package internal

import "time"

// Trip represents a WanderPlan travel itinerary.
type Trip struct {
	ID            string     `json:"id"`
	UserID        string     `json:"user_id"`
	Title         string     `json:"title"`
	Destination   string     `json:"destination"`
	CoverImageURL string     `json:"cover_image_url"`
	StartDate     *time.Time `json:"start_date"`
	EndDate       *time.Time `json:"end_date"`
	Status        string     `json:"status"`
	Visibility    string     `json:"visibility"`
	BudgetTotal   float64    `json:"budget_total"`
	Currency      string     `json:"currency"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ItineraryDay represents a single day in a trip.
type ItineraryDay struct {
	ID        string     `json:"id"`
	TripID    string     `json:"trip_id"`
	DayNumber int        `json:"day_number"`
	Date      *time.Time `json:"date"`
	Notes     string     `json:"notes"`
}

// ItineraryItem represents an activity, transport, accommodation or food entry.
type ItineraryItem struct {
	ID          string    `json:"id"`
	DayID       string    `json:"day_id"`
	TripID      string    `json:"trip_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	PlaceID     string    `json:"place_id"`
	Type        string    `json:"type"` // activity|transport|accommodation|food
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Cost        float64   `json:"cost"`
	Currency    string    `json:"currency"`
	OrderIndex  int       `json:"order_index"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTripRequest is the payload for POST /trips.
type CreateTripRequest struct {
	Title         string  `json:"title"         binding:"required"`
	Destination   string  `json:"destination"   binding:"required"`
	CoverImageURL string  `json:"cover_image_url"`
	StartDate     *string `json:"start_date"`
	EndDate       *string `json:"end_date"`
	Status        string  `json:"status"`
	Visibility    string  `json:"visibility"`
	BudgetTotal   float64 `json:"budget_total"`
	Currency      string  `json:"currency"`
}

// UpdateTripRequest is the payload for PATCH /trips/:id.
type UpdateTripRequest struct {
	Title         *string  `json:"title"`
	Destination   *string  `json:"destination"`
	CoverImageURL *string  `json:"cover_image_url"`
	StartDate     *string  `json:"start_date"`
	EndDate       *string  `json:"end_date"`
	Status        *string  `json:"status"`
	Visibility    *string  `json:"visibility"`
	BudgetTotal   *float64 `json:"budget_total"`
	Currency      *string  `json:"currency"`
}

// CreateDayRequest is the payload for POST /trips/:id/days.
type CreateDayRequest struct {
	DayNumber int     `json:"day_number" binding:"required"`
	Date      *string `json:"date"`
	Notes     string  `json:"notes"`
}

// UpdateDayRequest is the payload for PATCH /trips/:id/days/:dayId.
type UpdateDayRequest struct {
	Notes *string `json:"notes"`
	Date  *string `json:"date"`
}

// CreateItemRequest is the payload for POST /trips/:id/days/:dayId/items.
type CreateItemRequest struct {
	Title       string  `json:"title"    binding:"required"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	PlaceID     string  `json:"place_id"`
	Type        string  `json:"type"     binding:"required"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Cost        float64 `json:"cost"`
	Currency    string  `json:"currency"`
	OrderIndex  int     `json:"order_index"`
}

// UpdateItemRequest is the payload for PATCH /trips/:id/items/:itemId.
type UpdateItemRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Location    *string  `json:"location"`
	PlaceID     *string  `json:"place_id"`
	Type        *string  `json:"type"`
	StartTime   *string  `json:"start_time"`
	EndTime     *string  `json:"end_time"`
	Cost        *float64 `json:"cost"`
	Currency    *string  `json:"currency"`
	OrderIndex  *int     `json:"order_index"`
}

// ReorderRequest is the payload for PATCH /trips/:id/items/reorder.
type ReorderRequest struct {
	Items []ReorderItem `json:"items" binding:"required"`
}

// ReorderItem carries the new order_index for a single item.
type ReorderItem struct {
	ID         string `json:"id"          binding:"required"`
	OrderIndex int    `json:"order_index"`
}

// TripStats holds aggregate trip statistics for the recharts dashboard.
type TripStats struct {
	TotalTrips     int     `json:"total_trips"`
	TotalCountries int     `json:"total_countries"`
	TotalDays      int     `json:"total_days"`
	TotalBudget    float64 `json:"total_budget"`
}
