package models

import (
	"context"

	"github.com/projecteru2/core-plugins/volume/types"
)

// Realloc .
func (v *Volume) Realloc(ctx context.Context, node string, originResourceArgs *types.WorkloadResourceArgs, resourceOpts *types.WorkloadResourceOpts) (*types.EngineArgs, *types.WorkloadResourceArgs, *types.WorkloadResourceArgs, error) {
	return nil, nil, nil, nil
}
