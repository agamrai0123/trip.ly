---
applyTo: "backend/**/*_test.go"
---

# Go Test Rules

## Structure
- All tests are table-driven. Define a `tests []struct{ name, input, expected }` slice and range over it.
- Use `t.Run(tt.name, ...)` for every case so failures are clearly identified.
- Use `testify/require` for fatal assertions (stop the test on failure) and `testify/assert` for non-fatal ones.

## HTTP handler tests
- Use `httptest.NewRecorder()` and `httptest.NewRequest()` — never spin up a real server in unit tests.
- Build the Gin engine with `gin.New()` + only the middlewares under test; do not use `gin.Default()`.
- Assert the full response envelope shape, not just the status code.

## Integration tests
- Use `testcontainers-go` to start real PostgreSQL 16 and Kafka containers.
- Run migrations with `golang-migrate` against the container before each test suite.
- Tag integration tests with `//go:build integration` and run them separately from unit tests.

## Mocking
- Mock all gRPC clients with `mockery`-generated mocks. Never create hand-rolled mocks.
- Mock Kafka producers/consumers using the interfaces defined in `pkg/kafka`.
- Do not mock the DB in handler tests — use the real containerised DB.

## Coverage
- Business logic packages must reach ≥ 80% coverage. Run `go test -coverprofile=coverage.out ./...` to verify.
- Do not write trivial tests just to hit coverage. Every test case must assert meaningful behaviour.

## Naming
- Test files: `<file>_test.go` in the same package as the code under test.
- Test functions: `TestHandlerName_Scenario` (e.g. `TestCreateTrip_MissingTitle`).
- Helper functions: prefix with `must` (e.g. `mustCreateUser(t, db)`) and call `t.Helper()` inside them.
