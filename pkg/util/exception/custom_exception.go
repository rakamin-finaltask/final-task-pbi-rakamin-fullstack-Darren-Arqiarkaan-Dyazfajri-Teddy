package exception

import "errors"

var (
	NotFoundErr  = errors.New("not found")
	DuplicateErr = errors.New("value is duplicated")
)
