# ctx:search-service | 2026-05-24 | HTTP:8086
## state
todo: -
planned: Proto search.proto(RPCs:SearchPlaces,SearchTrips; comment all; run make proto) | scaffold+implement(go.mod,cmd,internal/*,air.toml,README.md; go work use) | Dockerfile(multi-stage,EXPOSE 8086,HEALTHCHECK /healthz,LABEL service=search-service) | Go tests(≥80%;mock Google Places API)
done: context-doc@2026-05-24
errors: -
_update when: Google Places integration/caching/FTS query/endpoints/new search types change_

## files
cmd/main.go: bootstrap — DB pool,Google Places client,gin,shutdown(SIGTERM/SIGINT 10s)
config/search-service-config.toml: TOML defaults
internal/config.go: Viper; DB DSN,GOOGLE_PLACES_API_KEY
internal/models.go: PlaceResult,TripSearchResult
internal/database.go: PlacesCacheRepo(GET/SET by cache_key); TripSearchRepo(tsvector FTS query)
internal/places.go: GooglePlacesClient: SearchNearby,SearchByText; check cache first; write on miss; rate-limit via x-goog-quota-user; NEVER log API key
internal/service.go: SearchService: SearchPlaces(cache-aside); SearchTrips(FTS+user_id filter)
internal/handlers.go: SearchPlacesHandler,SearchTripsHandler
internal/routes.go: RegisterRoutes; /search/* group
internal/metrics.go: search_requests_total,places_cache_hits_total,places_cache_misses_total,google_places_api_latency,db_query_duration; /metrics
internal/logger.go: zerolog+lumberjack | internal/errors.go: AppError helpers

## routes
GET /healthz,/metrics [no-auth]
GET /search/places?q=&lat=&lng= [JWT] — Google Places API + cache-aside
GET /search/trips?q= [JWT] — tsvector FTS on authenticated user's trips

## grpc: none | kafka: none

## env: PORT=8086 DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD GOOGLE_PLACES_API_KEY PLACES_CACHE_TTL_HOURS=24
## pkg: config database errors logger middleware response
## db: places_cache(SELECT by cache_key; INSERT on miss; TTL=cached_at<NOW()-INTERVAL'24h') | trips(SELECT via tsvector FTS+user_id)
## cache-key: SHA-256(q+"|"+lat+"|"+lng) hex
## fts: trips.search_vector=tsvector(title||destination||description); GIN index added in migration 000002; query=plainto_tsquery; ORDER BY ts_rank DESC

