package domain

import "errors"

// ErrNotFound is returned when a short code has no stored URL.
var ErrNotFound = errors.New("url not found")
