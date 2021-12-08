package types

import "errors"

var (
	ErrInvalidCapacity      = errors.New("invalid capacity")
	ErrInvalidVolume        = errors.New("invalid volume")
	ErrInsufficientResource = errors.New("cannot alloc a plan, not enough resource")

	ErrNodeExists = errors.New("node already exists")
)
