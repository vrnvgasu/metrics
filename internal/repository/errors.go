package repository

import "errors"

var (
	ErrNotFound   = errors.New("not found")
	ErrNotSupport = errors.New("not support")
)
