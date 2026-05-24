package internal

import pkgerr "github.com/agamrai0123/wanderplan/pkg/errors"

var (
	BadRequest   = pkgerr.BadRequest
	Unauthorized = pkgerr.Unauthorized
	Forbidden    = pkgerr.Forbidden
	NotFound     = pkgerr.NotFound
	Internal     = pkgerr.Internal
)
