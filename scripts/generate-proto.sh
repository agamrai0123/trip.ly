#!/usr/bin/env bash
# scripts/generate-proto.sh — Regenerate Go code from .proto files using buf.
#
# Equivalent to: make proto
# Requires: buf (install via `go install github.com/bufbuild/buf/cmd/buf@latest`)
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTO_DIR="$ROOT_DIR/backend/proto"

echo "▶  Generating proto code from $PROTO_DIR ..."

if ! command -v buf &>/dev/null; then
  echo "❌  buf not found. Install: go install github.com/bufbuild/buf/cmd/buf@latest"
  exit 1
fi

cd "$PROTO_DIR"
buf generate

echo "▶  Running go mod tidy in backend workspace..."
cd "$ROOT_DIR/backend"
go mod tidy 2>/dev/null || true  # workspace tidy; individual service modules are separate

echo "✅  Proto generation complete."
echo "   Generated files: $PROTO_DIR/gen/"
ls "$PROTO_DIR/gen/wanderplan/v1/" 2>/dev/null | head -20 || true
