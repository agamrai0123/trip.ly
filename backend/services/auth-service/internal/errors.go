package internal

// Re-export pkg/errors constructors for convenience.

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

var (
	BadRequest   = pkgerr.BadRequest
	Unauthorized = pkgerr.Unauthorized
	Internal     = pkgerr.Internal
)
