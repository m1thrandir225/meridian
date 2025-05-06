package persistence

import "errors"

var (
	ErrNotFound    = errors.New("requested item not found")
	ErrConcurrency = errors.New("database concurrency conflict")
)
