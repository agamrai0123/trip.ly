package internal

import "context"

// SearchService orchestrates place and trip search.
type SearchService struct {
	places   *PlacesClient
	tripRepo *TripSearchRepo
}

// NewSearchService constructs the service.
func NewSearchService(places *PlacesClient, tripRepo *TripSearchRepo) *SearchService {
	return &SearchService{places: places, tripRepo: tripRepo}
}

// SearchPlaces returns place results from Google Places (with caching).
func (s *SearchService) SearchPlaces(ctx context.Context, query string, lat, lng float64) ([]*PlaceResult, error) {
	return s.places.Search(ctx, query, lat, lng)
}

// SearchTrips performs full-text search over trips.
func (s *SearchService) SearchTrips(ctx context.Context, query string) ([]*TripSearchResult, error) {
	return s.tripRepo.Search(ctx, query)
}
