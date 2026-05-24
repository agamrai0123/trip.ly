---
applyTo: "backend/services/**/*.go"
---

# Go Service Handler Rules

## Error handling
- Never use `fmt.Errorf` or `errors.New` directly in handlers — always return `pkg/errors.AppError{Code, Message, HTTPStatus}`.
- Wrap every external call (DB, gRPC, Kafka) in an `AppError`; include the original error as the `Cause` field.
- Use `c.Error(err)` to attach errors to the Gin context; let the centralised error middleware write the response.

## Response shape
- Every successful response must use the standard envelope helper from `pkg/errors`:
  ```go
  c.JSON(http.StatusOK, envelope.OK(data, c.GetString("request_id")))
  ```
- Never call `c.JSON` with a raw struct — always go through the envelope.

## Database
- All queries use pgx named args (`@param`) or positional (`$1, $2`). No string concatenation in SQL, ever.
- Always pass `ctx` from the Gin context (`c.Request.Context()`) to every pgx call.
- Use `pgx.ErrNoRows` to detect not-found; map it to `AppError` with `HTTPStatus: 404`.

## Logging
- Use `zerolog` from `pkg/logger`. Attach `request_id` and `user_id` to every log entry.
- Log level: `Error` for unexpected failures, `Warn` for client errors (4xx), `Info` for normal operations.
- Never use `fmt.Print*` or the stdlib `log` package.

## Metrics
- Call the Prometheus counter/histogram helpers from `pkg/middleware` at the top of every handler.
- Instrument DB query duration using the helper in `pkg/database`. Do not create ad-hoc `prometheus.NewCounter` calls.

## gRPC clients
- Obtain gRPC clients from `pkg/grpc` factory — never dial directly inside a handler.
- Always set a deadline on outbound gRPC calls: `ctx, cancel := context.WithTimeout(ctx, 3*time.Second)`.

## Kafka
- Produce events via the `pkg/kafka` producer wrapper. Never import Sarama directly in service code.
- Topic names must come from constants defined in `pkg/kafka/topics.go`.

## Graceful shutdown
- Every `main.go` must register `SIGTERM`/`SIGINT` handlers and call the shared `pkg/shutdown.Drain(10s)` helper.
