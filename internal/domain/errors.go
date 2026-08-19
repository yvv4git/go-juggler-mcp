package domain

import "errors"

// ErrBrowserUnavailable is returned when the camofox-browser endpoint
// cannot be reached.
var ErrBrowserUnavailable = errors.New("browser unavailable")

// ErrTabNotFound is returned when the requested tab does not exist.
var ErrTabNotFound = errors.New("tab not found")

// ErrSessionNotFound is returned when the requested session does not exist.
var ErrSessionNotFound = errors.New("session not found")
