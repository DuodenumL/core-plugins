package types

import "errors"

var (
	ErrInvalidCapacity      = errors.New("invalid capacity")
	ErrInsufficientResource = errors.New("cannot alloc a plan, not enough resource")
	ErrInvalidStorage       = errors.New("invalid storage")
	ErrInvalidVolume        = errors.New("invalid volume")

	ErrNodeExists = errors.New("node already exists")
)
