package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/storage/types"
)

// GetNodesDeployCapacity .
func (s *Storage) GetNodesDeployCapacity(ctx context.Context, nodes []string, opts *types.WorkloadResourceOpts) (map[string]*types.NodeCapacityInfo, int, error) {
	if err := opts.Validate(); err != nil {
		logrus.Errorf("[GetNodesDeployCapacity] invalid resource opts %+v, err: %v", opts, err)
		return nil, 0, err
	}

	capacityInfoMap := map[string]*types.NodeCapacityInfo{}
	total := 0
	for _, node := range nodes {
		resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
		if err != nil {
			logrus.Errorf("[GetNodesDeployCapacity] failed to get resource info of node %v, err: %v", node, err)
			return nil, 0, err
		}
		capacityInfo := s.doGetNodeCapacityInfo(node, resourceInfo, opts)
		if capacityInfo.Capacity > 0 {
			capacityInfoMap[node] = capacityInfo
			total += capacityInfo.Capacity
		}
	}

	return capacityInfoMap, total, nil
}

func (s *Storage) doGetNodeCapacityInfo(node string, resourceInfo *types.NodeResourceInfo, opts *types.WorkloadResourceOpts) *types.NodeCapacityInfo {
	capacityInfo := &types.NodeCapacityInfo{
		Node:   node,
		Weight: 1,
	}
	availableResourceInfo := resourceInfo.GetAvailableResource()

	if opts.StorageRequest == 0 {
		capacityInfo.Capacity = s.config.Scheduler.MaxDeployCount
	} else {
		capacityInfo.Capacity = int(availableResourceInfo.Storage / opts.StorageRequest)
	}
	if capacityInfo.Capacity > s.config.Scheduler.MaxDeployCount {
		capacityInfo.Capacity = s.config.Scheduler.MaxDeployCount
	}

	capacityInfo.Usage = float64(resourceInfo.Usage.Storage) / float64(resourceInfo.Capacity.Storage)
	capacityInfo.Rate = float64(opts.StorageRequest) / float64(resourceInfo.Capacity.Storage)

	return capacityInfo
}
