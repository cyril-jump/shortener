package errs

import "errors"

// Internal errors
var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrWasDeleted    = errors.New("was deleted")
	ErrNoContent     = errors.New("no contents for this users")
	ErrNetNotTrusted = errors.New("network is not trusted")
)
