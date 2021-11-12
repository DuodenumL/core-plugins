package models

import (
	"context"
	"github.com/projecteru2/core-plugins/cpumem/schedule"
	"github.com/projecteru2/core-plugins/cpumem/types"
	"github.com/sirupsen/logrus"
)

// SelectAvailableNodes .
func (c *CPUMem) SelectAvailableNodes(ctx context.Context, nodes []string, opts *types.WorkloadResourceOpts) (map[string]*types.NodeCapacityInfo, error) {
	capacityInfoMap := map[string]*types.NodeCapacityInfo{}
	for _, node := range nodes {
		resourceInfo, err := c.doGetNodeResourceInfo(ctx, node)
		if err != nil {
			logrus.Errorf("[SelectAvailableNodes] failed to get resource info of node %v, err: %v", node, err)
			return nil, err
		}
		capacityInfo := c.doGetNodeCapacityInfo(node, resourceInfo, opts)
		if capacityInfo.Capacity > 0 {
			capacityInfoMap[node] = capacityInfo
		}
	}

	return capacityInfoMap, nil
}

func (c *CPUMem) doGetNodeCapacityInfo(node string, resourceInfo *types.NodeResourceInfo, opts *types.WorkloadResourceOpts) *types.NodeCapacityInfo {
	availableResourceArgs := resourceInfo.GetAvailableResource()

	capacityInfo := &types.NodeCapacityInfo{
		Node:     node,
		Weight:   1,
	}

	// if cpu-bind is not required, then returns capacity by memory
	if !opts.CPUBind {
		// check if cpu is enough
		if opts.CPURequest > float64(len(resourceInfo.Capacity.CPUMap)) {
			return capacityInfo
		}

		// calculate by memory request
		capacityInfo.Capacity = int(availableResourceArgs.Memory / opts.MemRequest)
		capacityInfo.Usage = float64(resourceInfo.Usage.Memory) / float64(resourceInfo.Capacity.Memory)
		capacityInfo.Rate = float64(opts.MemRequest) / float64(resourceInfo.Capacity.Memory)

		return capacityInfo
	}

	// if cpu-bind is required, then returns capacity by cpu scheduling
	cpuPlans := schedule.GetCPUPlans(resourceInfo, nil, c.config.Scheduler.ShareBase, c.config.Scheduler.MaxShare, opts)
	capacityInfo.Capacity = len(cpuPlans)
	capacityInfo.Usage = resourceInfo.Usage.CPU / resourceInfo.Capacity.CPU
	capacityInfo.Rate = opts.CPURequest / resourceInfo.Capacity.CPU
	capacityInfo.Weight = 100

	return nil
}
