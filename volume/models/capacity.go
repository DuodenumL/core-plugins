package models

import (
	"context"
	"math"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/volume/schedule"
	"github.com/projecteru2/core-plugins/volume/types"
)

// GetNodesCapacity .
func (v *Volume) GetNodesCapacity(ctx context.Context, nodes []string, opts *types.WorkloadResourceOpts) (map[string]*types.NodeCapacityInfo, int, error) {
	if err := opts.Validate(); err != nil {
		logrus.Errorf("[GetNodesCapacity] invalid resource opts %+v, err: %v", opts, err)
		return nil, 0, err
	}

	capacityInfoMap := map[string]*types.NodeCapacityInfo{}
	total := 0
	for _, node := range nodes {
		resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
		if err != nil {
			logrus.Errorf("[GetNodesCapacity] failed to get resource info of node %v, err: %v", node, err)
			return nil, 0, err
		}
		capacityInfo := v.doGetNodeCapacityInfo(node, resourceInfo, opts)
		if capacityInfo.Capacity > 0 {
			capacityInfoMap[node] = capacityInfo
			if total == math.MaxInt || capacityInfo.Capacity == math.MaxInt {
				total = math.MaxInt
			} else {
				total += capacityInfo.Capacity
			}
		}
	}

	return capacityInfoMap, total, nil
}

func (v *Volume) doGetNodeCapacityInfo(node string, resourceInfo *types.NodeResourceInfo, opts *types.WorkloadResourceOpts) *types.NodeCapacityInfo {
	capacityInfo := &types.NodeCapacityInfo{
		Node:   node,
		Weight: 1,
	}

	volumePlans := schedule.GetVolumePlans(resourceInfo, opts.VolumesRequest, nil, v.config.Scheduler.MaxDeployCount)
	capacityInfo.Capacity = len(volumePlans)
	capacityInfo.Usage = float64(resourceInfo.Usage.Volumes.Total()) / float64(resourceInfo.Capacity.Volumes.Total())
	capacityInfo.Rate = float64(opts.VolumesRequest.TotalSize()) / float64(resourceInfo.Capacity.Volumes.Total())

	return capacityInfo
}
