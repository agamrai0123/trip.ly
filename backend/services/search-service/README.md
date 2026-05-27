# search-service

Place and trip search service for WanderPlan. Proxies Google Places API with caching and runs full-text search over trips.

## Purpose

- Proxy `GET /search/places` to Google Places API with a cache-aside layer (`places_cache` table)
- Run PostgreSQL `tsvector` full-text search over the authenticated user's trips
- Cache key: `SHA-256(q + "|" + lat + "|" + lng)` hex string
- Never logs the `GOOGLE_PLACES_API_KEY`

## Endpoints

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| `GET` | `/healthz` | None | Health check |
| `GET` | `/metrics` | None | Prometheus metrics |
| `GET` | `/search/places?q=&lat=&lng=` | JWT | Google Places search (cached) |
| `GET` | `/search/trips?q=` | JWT | Full-text search over user's trips |

### Places Response Shape

```json
{
  "data": [
    {
      "place_id": "ChIJ...",
      "name": "Eiffel Tower",
      "address": "Champ de Mars, Paris",
      "lat": 48.8584,
      "lng": 2.2945,
      "photos": ["https://..."]
    }
  ]
}
```

### Trip Search Response Shape

```json
{
  "data": [
    {
      "id": "uuid",
      "title": "Euro Summer 2026",
      "destination": "Paris, France",
      "start_date": "2026-07-01",
      "end_date": "2026-07-14"
    }
  ]
}
```

## Caching Strategy

1. Hash the query + coordinates to a cache key
2. `SELECT` from `places_cache` where `cache_key = $1 AND cached_at > NOW() - INTERVAL '24h'`
3. On hit: return cached result
4. On miss: call Google Places API, `INSERT` result, return

## Full-Text Search

- Column: `trips.search_vector` (`tsvector` of `title || destination || description`)
- Index: GIN (added in migration `000002`)
- Query: `plainto_tsquery($1)` ordered by `ts_rank DESC`, filtered by `user_id`

## Environment Variables

| Variable | Description |
|----------|-------------|
| `PORT` | HTTP port (default: `8086`) |
| `DB_HOST` | PostgreSQL host |
| `DB_PORT` | PostgreSQL port |
| `DB_NAME` | Database name |
| `DB_USER` | Database user |
| `DB_PASSWORD` | Database password |
| `DB_SSL_MODE` | `disable` / `require` |
| `GOOGLE_PLACES_API_KEY` | Google Places API key |
| `PLACES_CACHE_TTL_HOURS` | Cache TTL in hours (default: `24`) |
| `RATE_LIMIT_RPS` | Rate limit — requests per second |
| `RATE_LIMIT_BURST` | Rate limit — burst size |

## How to Run Locally

```bash
cp .env.example .env   # add GOOGLE_PLACES_API_KEY
cd backend/services/search-service
air
```

## How to Run Tests

```bash
cd backend/services/search-service
go test ./... -v -count=1

# Integration tests (requires Docker)
go test ./... -v -count=1 -tags integration
```

## Observability

- `GET /metrics` — `search_requests_total`, `places_cache_hits_total`, `places_cache_misses_total`, `google_places_api_latency_seconds`, `db_query_duration_seconds`
- Structured JSON logs via `zerolog` + `lumberjack`

## Database Tables

- `places_cache` — cache by SHA-256 key; TTL filter on `cached_at`
- `trips` — `tsvector` FTS query with `user_id` filter
