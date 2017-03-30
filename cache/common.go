package cache

import "errors"

var (
	// ErrorMissed code.
	ErrorMissed = errors.New("cache: missed")

	// ErrorPtr code.
	ErrorPtr = errors.New("cache: ptr target")

	// ErrorSerialize code.
	ErrorSerialize = errors.New("cache: failed to serialize data")

	// ErrorUnserialize code.
	ErrorUnserialize = errors.New("cache: failed to unserialize data")

	// ErrorTypeMissmatch code.
	ErrorTypeMissmatch = errors.New("cache: return type missmatch")
)
