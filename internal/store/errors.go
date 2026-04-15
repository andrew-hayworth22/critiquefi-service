package store

import "errors"

var (
	ErrNotFound            = errors.New("not found")
	ErrDuplicate           = errors.New("duplicate record")
	ErrForeignKeyViolation = errors.New("foreign key violation")
	ErrNotNullViolation    = errors.New("not null violation")
)
