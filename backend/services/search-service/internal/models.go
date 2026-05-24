package internal

// PlaceResult is a single result from the Google Places API.
type PlaceResult struct {
	PlaceID  string   `json:"place_id"`
	Name     string   `json:"name"`
	Address  string   `json:"formatted_address"`
	Types    []string `json:"types"`
	Lat      float64  `json:"lat"`
	Lng      float64  `json:"lng"`
	Rating   float64  `json:"rating,omitempty"`
	PhotoRef string   `json:"photo_reference,omitempty"`
}

// TripSearchResult is a trip returned by full-text search.
type TripSearchResult struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Destination string `json:"destination"`
	Status      string `json:"status"`
	OwnerID     string `json:"owner_id"`
}
