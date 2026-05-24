package internal

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

// Re-export pkg/errors constructors for convenience inside this package.
var (
	BadRequest   = pkgerr.BadRequest
	Unauthorized = pkgerr.Unauthorized
	Internal     = pkgerr.Internal
	NotFound     = pkgerr.NotFound
)
