# ctx:collaboration-service | 2026-05-24 | HTTP:8084
## state
todo: -
planned: Dockerfile(multi-stage,EXPOSE 8084,HEALTHCHECK /healthz)
done: context-doc@2026-05-24 | scaffold+implement@2026-05-24 | .env.example@2026-05-25 | README.md@2026-05-25 | tests(handlers_test.go+service_test.go)@2026-05-25 | proto(collaboration.proto+gen)@2026-05-25
errors: -
_update when: collaborator roles/endpoints/Kafka events/collaborators schema change_

## files
cmd/main.go: bootstrap — DB pool,Kafka producer,gin,shutdown(SIGTERM/SIGINT 10s)
config/collaboration-service-config.toml: TOML defaults
internal/config.go: Viper; DB DSN,Kafka brokers
internal/models.go: Collaborator; InviteRequest,UpdateRoleRequest; Role=viewer|editor|admin
internal/database.go: CollaboratorRepo: List,Invite(insert),UpdateRole,Remove; parameterised pgx; reads users+trips for validation
internal/service.go: CollaborationService; validates trip ownership; resolves invitee email→user_id(SELECT users); publishes Kafka event
internal/handlers.go: ListCollaborators,InviteCollaborator,UpdateCollaboratorRole,RemoveCollaborator
internal/routes.go: RegisterRoutes; /trips/:id/collaborators group
internal/metrics.go: collab_operations_total,db_query_duration,goroutine_count; /metrics
internal/logger.go: zerolog+lumberjack | internal/errors.go: AppError helpers

## routes
GET /healthz,/metrics [no-auth]
GET /trips/:id/collaborators [JWT]
POST /trips/:id/collaborators [JWT,owner-only] — invite by email
PATCH /trips/:id/collaborators/:userId [JWT,owner-only] — change role
DELETE /trips/:id/collaborators/:userId [JWT,owner-only] — remove

## grpc: none
## kafka-out: collab-events:collaboration.invited{trip_id,inviter_user_id,invitee_user_id,invitee_email,role}

## env: PORT=8084 DB_HOST DB_PORT DB_NAME DB_USER DB_PASSWORD KAFKA_BROKERS
## pkg: config database errors kafka logger middleware response
## db: collaborators(SELECT by trip_id,INSERT,UPDATE role,DELETE) | users(SELECT by email for invitee resolution) | trips(SELECT owner_id for ownership check)
## authz: only owner can invite/change-role/remove; cannot elevate own role; removing owner→400; hierarchy:owner>admin>editor>viewer

