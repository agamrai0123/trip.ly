package internal

// Re-export pkg/errors constructors for convenience.

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

// Convenience re-exports of pkg/errors constructors used throughout this package.
var (
	BadRequest   = pkgerr.BadRequest
	Unauthorized = pkgerr.Unauthorized
	Internal     = pkgerr.Internal
)
