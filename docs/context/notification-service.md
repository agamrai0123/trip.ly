# ctx:notification-service | 2026-05-24 | HTTP:8085 WS:8085
## state
todo: -
planned: -
done: context-doc@2026-05-24 | scaffold+implement@2026-05-24 | .env.example@2026-05-25 | README.md@2026-05-25 | tests(handlers_test.go+service_test.go)@2026-05-25 | notification.proto+gen@2026-05-25 | db_pool_acquired_connections metric@2026-05-25 | proto(notification.proto+gen)@2026-05-25 | Dockerfile(multi-stage,EXPOSE 8085,HEALTHCHECK,LABEL)@2026-05-25 | kafka_consumer_lag GaugeVec(topic,partition labels)in consumer.go;reg passed to NewEventConsumer@2026-05-26
errors: -
_update when: Kafka topics/event types/WebSocket logic/endpoints/notifications schema/consumer group change_

## files
cmd/main.go: bootstrap — DB pool,Kafka consumer group init,WebSocket Hub start,gin,shutdown(SIGTERM/SIGINT 10s)
config/notification-service-config.toml: TOML defaults
internal/config.go: Viper; DB DSN,Kafka brokers,consumer group name
internal/models.go: Notification; EventPayload structs for user.login,trip.created/updated/deleted,collaboration.invited
internal/database.go: NotificationRepo: List(paginated),Insert,MarkRead,MarkAllRead; parameterised pgx
internal/consumer.go: Sarama ConsumerGroupHandler; ConsumeClaim dispatches by event type; creates Notification+Hub.Broadcast
internal/websocket.go: Hub{clients map,broadcast chan,register/unregister chans}+Client{conn,send chan}; Hub.Run goroutine; ReadPump/WritePump
internal/service.go: NotificationService wraps repo+hub; GetNotifications,MarkRead,MarkAllRead
internal/handlers.go: ListHandler,MarkReadHandler,MarkAllReadHandler,WebSocketHandler(HTTP→WS upgrade)
internal/routes.go: RegisterRoutes; /notifications/* + /ws/notifications
internal/metrics.go: notifications_created_total,kafka_consumer_lag(REQUIRED GaugeVec{topic,partition}),websocket_active_connections,db_query_duration,goroutine_count; /metrics; updated after MarkOffset
internal/logger.go: zerolog+lumberjack | internal/errors.go: AppError helpers

## routes
GET /healthz,/metrics [no-auth]
GET /notifications [JWT] — paginated newest-first
PATCH /notifications/:id/read [JWT]
PATCH /notifications/read-all [JWT]
GET /ws/notifications [JWT via query-param or header] — WebSocket upgrade

## grpc: none
## kafka-in: auth-events(user.login→notif:new-device-login) | trip-events(trip.created/updated/deleted→notif for collaborators) | collab-events(collaboration.invited→notif for invitee); group:wanderplan-notifications; BalanceStrategyRange

## env: PORT=8085 DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD KAFKA_BROKERS KAFKA_CONSUMER_GROUP=wanderplan-notifications
## pkg: config database errors kafka logger middleware response
## db: notifications(INSERT from Kafka,SELECT paginated by user_id,UPDATE read_at)
## ws-arch: one Hub per process; one Client per tab; Hub.broadcast→client.send→WritePump→WS frame; msg={type:"notification",data:<Notification>}; Hub not persistent—DB is source of truth

