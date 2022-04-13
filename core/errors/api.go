package errors

import (
	"fmt"
)

// Common Service errors.
var (
	InternalServerError    = NewHTTPInternal(100000, "internal server error")
	UserAuthorizationError = NewHTTPUnauthorized(100001, "unauthorized request")
	EntityNotFoundError    = func(entityType string) *HTTP {
		return NewHTTPNotFound(100002, fmt.Sprintf(
			"`%s` entity not found", entityType),
		)
	}
	PathParametersParsingError = NewHTTPBadRequest(
		100003, "could not parse the path parameters",
	)
	QueryParametersParsingError = NewHTTPBadRequest(
		100004, "could not parse the query parameters",
	)
)
