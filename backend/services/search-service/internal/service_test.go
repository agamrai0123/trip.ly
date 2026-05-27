package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────────────────────────────────────
// PlaceResult — JSON field name contract
// ──────────────────────────────────────────────────────────────────────────────

// TestPlaceResult_JSONFields checks that PlaceResult uses the JSON field names
// expected by the frontend API contract.
func TestPlaceResult_JSONFields(t *testing.T) {
	p := PlaceResult{
		PlaceID:  "ChIJN1t_tDeuEmsRUsoyG83frY4",
		Name:     "Sydney Opera House",
		Address:  "Bennelong Point, Sydney NSW 2000, Australia",
		Types:    []string{"tourist_attraction", "point_of_interest"},
		Lat:      -33.8568,
		Lng:      151.2153,
		Rating:   4.7,
		PhotoRef: "photo-ref-123",
	}

	b, err := json.Marshal(p)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))

	assert.Equal(t, "ChIJN1t_tDeuEmsRUsoyG83frY4", m["place_id"])
	assert.Equal(t, "Sydney Opera House", m["name"])
	assert.Equal(t, "Bennelong Point, Sydney NSW 2000, Australia", m["formatted_address"])
	assert.InDelta(t, -33.8568, m["lat"], 0.0001)
	assert.InDelta(t, 151.2153, m["lng"], 0.0001)
	assert.InDelta(t, 4.7, m["rating"], 0.001)
	assert.Equal(t, "photo-ref-123", m["photo_reference"])
}

// TestPlaceResult_OptionalFieldsOmitted checks that zero-value optional fields
// are omitted from JSON (omitempty tags).
func TestPlaceResult_OptionalFieldsOmitted(t *testing.T) {
	p := PlaceResult{
		PlaceID: "place-001",
		Name:    "Central Park",
		Address: "New York, NY",
		Lat:     40.7851,
		Lng:     -73.9683,
	}

	b, err := json.Marshal(p)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))

	// Rating == 0 and PhotoRef == "" → should be omitted due to omitempty
	_, hasRating := m["rating"]
	_, hasPhoto := m["photo_reference"]
	assert.False(t, hasRating, "zero rating should be omitted")
	assert.False(t, hasPhoto, "empty photo_reference should be omitted")
}

// ──────────────────────────────────────────────────────────────────────────────
// TripSearchResult — JSON field name contract
// ──────────────────────────────────────────────────────────────────────────────

// TestTripSearchResult_JSONFields verifies the search result JSON shape.
func TestTripSearchResult_JSONFields(t *testing.T) {
	tests := []struct {
		name        string
		result      TripSearchResult
		wantTitle   string
		wantDestKey string
	}{
		{
			name: "draft trip",
			result: TripSearchResult{
				ID:          "t-001",
				Title:       "Tokyo Adventure",
				Destination: "Tokyo, Japan",
				Status:      "draft",
				OwnerID:     "u-001",
			},
			wantTitle:   "Tokyo Adventure",
			wantDestKey: "destination",
		},
		{
			name: "published trip",
			result: TripSearchResult{
				ID:          "t-002",
				Title:       "Paris Escape",
				Destination: "Paris, France",
				Status:      "published",
				OwnerID:     "u-002",
			},
			wantTitle:   "Paris Escape",
			wantDestKey: "destination",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.result)
			require.NoError(t, err)

			var m map[string]interface{}
			require.NoError(t, json.Unmarshal(b, &m))

			assert.Equal(t, tt.wantTitle, m["title"])
			assert.Contains(t, m, tt.wantDestKey)
			assert.Equal(t, tt.result.Destination, m["destination"])
			assert.Equal(t, tt.result.Status, m["status"])
			assert.Equal(t, tt.result.OwnerID, m["owner_id"])
		})
	}
}
