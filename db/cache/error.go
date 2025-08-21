package cache

import "errors"

var (
	ErrNotFound = errors.New("not found")

	ErrSetCache = errors.New("could not set cache")
)
