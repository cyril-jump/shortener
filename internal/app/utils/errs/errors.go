package errs

import "errors"

var (
	NotFound      = errors.New("not found")
	AlreadyExists = errors.New("already exists")
)
