package command

import "errors"

var (
	ErrNoSuchMethod  = errors.New("no such method")
	ErrInvalidParams = errors.New("invalid params")
)
