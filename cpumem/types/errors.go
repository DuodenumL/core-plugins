package types

import "errors"

var (
	ErrInvalidCapacity   = errors.New("invalid capacity")
	ErrInvalidUsage      = errors.New("invalid usage")
	ErrInvalidCPUMap     = errors.New("invalid cpu map")
	ErrInvalidNUMA       = errors.New("invalid numa")
	ErrInvalidNUMAMemory = errors.New("invalid numa")
	ErrInvalidMemory     = errors.New("invalid memory")
	ErrInvalidCPU        = errors.New("invalid cpu")

	ErrInsufficientCPU      = errors.New("cannot alloc a plan, not enough cpu")
	ErrInsufficientMEM      = errors.New("cannot alloc a plan, not enough memory")
	ErrInsufficientResource = errors.New("cannot alloc a plan, not enough resource")

	ErrNodeExists = errors.New("node already exists")
)
