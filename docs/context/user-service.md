# ctx:user-service | 2026-05-24 | HTTP:8083 gRPC:9083
## state
todo: scaffold+implement — go.mod,cmd/main.go,internal/*(config,handlers,routes,service,models,metrics,logger,errors,database),air.toml,README.md; `go work use ./services/user-service` in backend/
planned: Dockerfile(multi-stage,EXPOSE 8083 9083,HEALTHCHECK /healthz,LABEL service=user-service) | Go tests(≥80%)
done: context-doc@2026-05-24
errors: -
_update when: user profile fields/endpoints/gRPC clients/DB schema change_

## files
cmd/main.go: bootstrap — DB pool,gRPC client→trip-service,gin,shutdown(SIGTERM/SIGINT 10s)
config/user-service-config.toml: TOML defaults
internal/config.go: Viper; DB DSN,TRIP_SERVICE_GRPC_ADDR
internal/models.go: User,UpdateUserRequest,StatsResponse
internal/database.go: UserRepo: GetByID,UpdateProfile; parameterised pgx
internal/service.go: UserService.GetMe,UpdateMe; GetMyTrips→gRPC TripService.ListTripsByUser; GetMyStats→gRPC TripService.GetTripStats
internal/handlers.go: GetMeHandler,UpdateMeHandler,GetMyTripsHandler,GetMyStatsHandler
internal/routes.go: RegisterRoutes; /users/* group
internal/metrics.go: user_operations_total,grpc_call_duration,db_query_duration,goroutine_count; /metrics
internal/logger.go: zerolog+lumberjack | internal/errors.go: AppError helpers

## routes
GET /healthz,/metrics [no-auth]
GET /users/me [JWT] | PATCH /users/me [JWT] — update name,avatar_url
GET /users/me/trips [JWT] — delegates→TripService.ListTripsByUser gRPC
GET /users/me/stats [JWT] — delegates→TripService.GetTripStats gRPC → recharts

## grpc-out: trip-service TripService.ListTripsByUser+GetTripStats @TRIP_SERVICE_GRPC_ADDR
## grpc-in: UserService.GetUser(user_id)→{id,email,name,avatar_url}; not called yet (reserved)
## kafka: none

## env: PORT=8083 GRPC_PORT=9083 DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD TRIP_SERVICE_GRPC_ADDR=localhost:9082
## pkg: config database errors grpc logger middleware response
## db: users(SELECT by ID for GetMe; UPDATE name/avatar_url)

