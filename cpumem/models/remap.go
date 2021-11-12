package models

import (
	"context"
	"github.com/projecteru2/core-plugins/cpumem/types"
	"github.com/sirupsen/logrus"
)

// Remap .
func (c *CPUMem) Remap(ctx context.Context, node string, workloadResourceMap map[string]*types.WorkloadResourceArgs) (map[string]*types.EngineArgs, error) {
	resourceInfo, err := c.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[Remap] failed to get resource info of node %v, err: %v", node, err)
		return nil, err
	}
	availableNodeResource := resourceInfo.GetAvailableResource()

	shareCPUMap := types.CPUMap{}
	for cpu, pieces := range availableNodeResource.CPUMap {
		if pieces >= c.config.Scheduler.ShareBase {
			shareCPUMap[cpu] = c.config.Scheduler.ShareBase
		}
	}

	engineArgsMap := map[string]*types.EngineArgs{}

	for workloadID, workloadResourceArgs := range workloadResourceMap {
		// only process workloads without cpu limit
		if len(workloadResourceArgs.CPUMap) == 0 && workloadResourceArgs.CPULimit == 0 {
			engineArgsMap[workloadID] = &types.EngineArgs{
				CPU:      0,
				CPUMap:   shareCPUMap,
				NUMANode: "",
				Memory:   workloadResourceArgs.MemoryLimit,
			}
		}
	}
	return engineArgsMap, nil
}
