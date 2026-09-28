package repository

import "errors"

// ErrNotFound is returned by every GetXByID method when the requested
// document does not exist. Callers should compare with errors.Is.
var ErrNotFound = errors.New("document not found")
