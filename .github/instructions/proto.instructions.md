---
applyTo: "backend/proto/**/*.proto"
---

# Protobuf / gRPC Rules

## File layout
- One service per `.proto` file. File name matches the service: `auth.proto`, `trip.proto`, etc.
- Package name: `wanderplan.<service>` (e.g. `package wanderplan.auth`).
- Go package option: `option go_package = "github.com/wanderplan/backend/proto/<service>";`

## Comments
- Every `service`, `rpc`, `message`, and `field` must have a comment explaining its purpose.
- Use `//` line comments directly above the element — not `/* */` block comments.

## Naming conventions
- RPC methods: PascalCase verb-noun (e.g. `ValidateToken`, `GetTrip`, `CreateCollaborator`).
- Message names: PascalCase noun (e.g. `TokenRequest`, `UserClaims`).
- Field names: snake_case (e.g. `user_id`, `expires_at`).

## Versioning
- Do not remove or renumber existing fields. Add new fields only with new field numbers.
- Mark deprecated fields with `[deprecated = true]` and a comment explaining the replacement.

## Error handling
- Use standard gRPC status codes. Map `AppError.HTTPStatus` to the correct `codes.Code` in the service adapter.

## Codegen
- Always regenerate Go code with `make proto` after any `.proto` change. Never edit generated `*_grpc.pb.go` or `*.pb.go` files manually.
