package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/storage/types"
)

// Realloc .
func (s *Storage) Realloc(ctx context.Context, node string, originResourceArgs *types.WorkloadResourceArgs, resourceOpts *types.WorkloadResourceOpts) (*types.EngineArgs, *types.WorkloadResourceArgs, *types.WorkloadResourceArgs, error) {
	resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[GetReallocArgs] failed to get resource info of node %v, err: %v", node, err)
		return nil, nil, nil, err
	}
	availableResourceInfo := resourceInfo.GetAvailableResource()

	finalWorkloadResourceArgs := &types.WorkloadResourceArgs{
		StorageRequest: resourceOpts.StorageRequest,
		StorageLimit:   resourceOpts.StorageLimit,
	}

	if finalWorkloadResourceArgs.StorageRequest > availableResourceInfo.Storage {
		return nil, nil, nil, types.ErrInsufficientResource
	}
	engineArgs := &types.EngineArgs{Storage: finalWorkloadResourceArgs.StorageLimit}
	deltaWorkloadResourceArgs := &types.WorkloadResourceArgs{StorageRequest: finalWorkloadResourceArgs.StorageRequest - originResourceArgs.StorageRequest}

	return engineArgs, deltaWorkloadResourceArgs, finalWorkloadResourceArgs, nil
}
