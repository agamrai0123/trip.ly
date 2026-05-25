# ctx:shared-packages | 2026-05-24 | backend/pkg/
## state
todo: -
planned: -
done: all 7 packages implemented@2026-05-24 | context-doc@2026-05-24 | Migration 000001_init(fixed:user_id,budget_total,visibility,trip_id on items,order_index)@2025-07-26 | Migration 000003_fix_trips_and_items(ALTER for live DB that had old 000001)@2025-07-26 | Migration 000002_add_search_vector(search_vector+GIN on trips;GIN on places_cache.results)@2026-05-25 | pkg/database/metrics.go(PgxPoolCollector)@2026-05-25 | pkg/middleware.Metrics():removed namespace/subsystem prefix(metrics now http_requests_total+http_request_duration_seconds unprefixed);fixed status label to numeric strconv.Itoa();removed broken goroutines gauge@2026-05-26
errors: -
_update when: package added/function sig changed/new config keys/error codes/new middleware_

## pkg/config
Load(path,target) — reads TOML; overrides any field with matching env var; unmarshals into target; used-by: all 7 services

## pkg/database
Connect(ctx,dsn)→*pgxpool.Pool — MaxConns=20,MinConns=2,MaxConnLifetime=30min; pings on connect; used-by: auth,trip,user,collab,notification,search

## pkg/errors
AppError{Code,Message,Detail,HTTPStatus}; constructors: NotFound,InvalidInput,Unauthorized,Forbidden,Internal,Conflict
handlers call response.Err(c,appErr); used-by: all

## pkg/grpc
NewClient(addr)→*grpc.ClientConn — insecure,5s timeout+retry; used-by: api-gateway(→auth),user-service(→trip)

## pkg/jwt
Manager{Sign(Claims)→RS256-JWT(15min-exp), Verify(token)→*Claims}; NewManagerFromBase64(privB64,pubB64)
Claims{UserID,Email,Name,AvatarURL,RegisteredClaims}; used-by: auth-service(sign+verify), api-gateway(verify via AuthValidator gRPC)

## pkg/kafka
Producer{Publish(topic,Event)}; Consumer{Subscribe(ctx,topics,handler)}; Close() on each
Event{Type string,Timestamp time.Time,Payload map[string]any}
Topics: TopicAuthEvents="auth-events",TopicTripEvents="trip-events",TopicCollabEvents="collab-events"
Producers: auth,trip,collab; Consumer: notification-service only

## pkg/logger
InitLogger(level,logPath)→zerolog.Logger — stdout(JSON)+lumberjack(100MB,7days,5backups); used-by: all

## pkg/middleware
RequestID(): UUID v4→X-Request-ID+gin ctx
Logger(log): method,path,status,latency,request_id via zerolog
CORS(origins): headers+credentials+OPTIONS preflight
Auth(validator): Bearer token→validator.ValidateToken(TokenValidator interface)→inject "user" claims; no concrete type (avoids circular imports)
Recovery(): panic→500 AppError envelope
RateLimit(rps,burst): per-IP token bucket→429
Metrics(reg): request_count counter+latency histogram per req

## pkg/response
Envelope{Data,Error,Meta{RequestID,Timestamp(RFC3339)}}
OK(200),Created(201),NoContent(204),Err(appErr.HTTPStatus),ErrRaw(status,msg); used-by: all

## rules: pkg/* must NOT import services/*; services import pkg/* not each other (use gRPC/Kafka)

