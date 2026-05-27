package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const googlePlacesURL = "https://maps.googleapis.com/maps/api/place/textsearch/json"

// PlacesClient wraps the Google Places API with a local cache.
type PlacesClient struct {
	apiKey     string
	httpClient *http.Client
	cache      *PlaceCacheRepo
}

// NewPlacesClient creates a new client.
func NewPlacesClient(apiKey string, cache *PlaceCacheRepo) *PlacesClient {
	return &PlacesClient{
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
		cache:      cache,
	}
}

// Search returns places matching query, using the cache when available.
func (c *PlacesClient) Search(ctx context.Context, query string, lat, lng float64) ([]*PlaceResult, error) {
	cacheKey := fmt.Sprintf("%s|%.4f|%.4f", query, lat, lng)

	if c.cache != nil {
		if cached, err := c.cache.Get(ctx, cacheKey); err == nil {
			return cached, nil
		}
	}

	if c.apiKey == "" {
		return []*PlaceResult{}, nil
	}

	params := url.Values{}
	params.Set("query", query)
	params.Set("key", c.apiKey)
	if lat != 0 || lng != 0 {
		params.Set("location", fmt.Sprintf("%f,%f", lat, lng))
		params.Set("radius", "50000")
	}

	reqURL := googlePlacesURL + "?" + params.Encode()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var apiResp struct {
		Results []struct {
			PlaceID          string   `json:"place_id"`
			Name             string   `json:"name"`
			FormattedAddress string   `json:"formatted_address"`
			Types            []string `json:"types"`
			Geometry         struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
			} `json:"geometry"`
			Rating float64 `json:"rating"`
			Photos []struct {
				PhotoReference string `json:"photo_reference"`
			} `json:"photos"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, err
	}

	places := make([]*PlaceResult, 0, len(apiResp.Results))
	for _, r := range apiResp.Results {
		p := &PlaceResult{
			PlaceID: r.PlaceID,
			Name:    r.Name,
			Address: r.FormattedAddress,
			Types:   r.Types,
			Lat:     r.Geometry.Location.Lat,
			Lng:     r.Geometry.Location.Lng,
			Rating:  r.Rating,
		}
		if len(r.Photos) > 0 {
			p.PhotoRef = r.Photos[0].PhotoReference
		}
		places = append(places, p)
	}

	if c.cache != nil {
		_ = c.cache.Set(ctx, cacheKey, places)
	}
	return places, nil
}

// marshalPlaces serialises place results to JSON bytes.
func marshalPlaces(places []*PlaceResult) ([]byte, error) {
	return json.Marshal(places)
}

// unmarshalPlaces deserialises place results from JSON bytes.
func unmarshalPlaces(data []byte) ([]*PlaceResult, error) {
	var places []*PlaceResult
	if err := json.Unmarshal(data, &places); err != nil {
		return nil, err
	}
	return places, nil
}
