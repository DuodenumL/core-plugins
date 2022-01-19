package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/storage/types"
)

// GetDeployArgs .
func (s *Storage) GetDeployArgs(ctx context.Context, node string, deployCount int, opts *types.WorkloadResourceOpts) ([]*types.EngineArgs, []*types.WorkloadResourceArgs, error) {
	if err := opts.Validate(); err != nil {
		logrus.Errorf("[GetDeployArgs] invalid resource opts %+v, err: %v", opts, err)
		return nil, nil, err
	}

	resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[GetDeployArgs] failed to get resource info of node %v, err: %v", node, err)
		return nil, nil, err
	}

	return s.doAlloc(resourceInfo, deployCount, opts)
}

func (s *Storage) doAlloc(resourceInfo *types.NodeResourceInfo, deployCount int, opts *types.WorkloadResourceOpts) ([]*types.EngineArgs, []*types.WorkloadResourceArgs, error) {
	availableResourceInfo := resourceInfo.GetAvailableResource()
	capacity := int(availableResourceInfo.Storage / opts.StorageRequest)
	if capacity > s.config.Scheduler.MaxDeployCount {
		capacity = s.config.Scheduler.MaxDeployCount
	}

	if capacity < deployCount {
		return nil, nil, types.ErrInsufficientResource
	}

	resEngineArgs := []*types.EngineArgs{}
	resResourceArgs := []*types.WorkloadResourceArgs{}

	for i := 0; i < deployCount; i++ {
		resEngineArgs = append(resEngineArgs, &types.EngineArgs{Storage: opts.StorageLimit})
		resResourceArgs = append(resResourceArgs, &types.WorkloadResourceArgs{StorageRequest: opts.StorageRequest, StorageLimit: opts.StorageLimit})
	}

	return resEngineArgs, resResourceArgs, nil
}
