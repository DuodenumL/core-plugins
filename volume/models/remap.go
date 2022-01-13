package models

import (
	"context"

	"github.com/projecteru2/core-plugins/volume/types"
)

// Remap .
func (v *Volume) Remap(ctx context.Context, node string, workloadResourceMap *types.WorkloadResourceArgsMap) (map[string]*types.EngineArgs, error) {
	return nil, nil
}
