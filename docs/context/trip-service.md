# ctx:trip-service | 2026-05-24 | HTTP:8082 gRPC:9082
## state
todo: -
planned: Dockerfile(multi-stage,EXPOSE 8082 9082,HEALTHCHECK /healthz,LABEL service=trip-service)
done: context-doc@2026-05-24 | scaffold+implement@2026-05-24 | .env.example@2026-05-25 | README.md@2026-05-25 | tests(service_test.go+handlers_test.go+repository_test.go)@2026-05-25
errors: -
_update when: endpoints/DB/Kafka/gRPC/reorder/duplicate logic change_

## files
cmd/main.go: bootstrap — DB pool,Kafka producer,gRPC server,gin,shutdown(SIGTERM/SIGINT 10s)
config/trip-service-config.toml: TOML defaults
internal/config.go: Viper; DB DSN,Kafka brokers
internal/models.go: Trip,ItineraryDay,ItineraryItem; CreateTripReq,UpdateTripReq; ReorderRequest([]{id,order_index})
internal/database.go: TripRepo+DayRepo+ItemRepo; parameterised pgx; transaction for bulk reorder
internal/service.go: TripService; DuplicateTrip(deep-copy days+items,new IDs,no collaborators); ReorderItems(bulk UPDATE in tx,full-rollback on fail); publish Kafka after DB writes
internal/handlers.go: ListTrips,CreateTrip,GetTrip,UpdateTrip,DeleteTrip,DuplicateTrip; day CRUD; item CRUD; ReorderItems
internal/routes.go: RegisterRoutes; /trips/* groups
internal/helpers.go: extractUserID,validateOwnership(trips.owner_id==JWT.user_id else 403),buildTripResponse
internal/metrics.go: trip_operations_total,db_query_duration,goroutine_count; /metrics
internal/logger.go: zerolog+lumberjack | internal/errors.go: AppError helpers

## routes
GET /healthz,/metrics [no-auth]
GET,POST /trips [JWT] | GET,PATCH,DELETE /trips/:id [JWT] | POST /trips/:id/duplicate [JWT]
GET,POST /trips/:id/days [JWT] | PATCH,DELETE /trips/:id/days/:dayId [JWT]
GET,POST /trips/:id/days/:dayId/items [JWT] | PATCH,DELETE /trips/:id/items/:itemId [JWT]
PATCH /trips/:id/items/reorder [JWT] — body:[{id,order_index}] bulk update (dnd-kit)

## grpc-in: TripService.ListTripsByUser(user_id)→[TripSummary] | TripService.GetTripStats(user_id)→{total_trips,total_countries,total_days,total_budget}; called-by:user-service
## kafka-out: trip-events:trip.created{trip_id,user_id,title,destination} | trip.updated{trip_id,user_id,changed_fields[]} | trip.deleted{trip_id,user_id}

## env: PORT=8082 GRPC_PORT=9082 DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD KAFKA_BROKERS
## pkg: config database errors kafka logger middleware response
## db: trips(INSERT,SELECT,UPDATE,DELETE) | itinerary_days(ordered by day_number) | itinerary_items(bulk UPDATE order_index in tx) | trip_tags(batch INSERT,SELECT,DELETE)

