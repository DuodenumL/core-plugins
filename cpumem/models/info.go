package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/cpumem/types"
	"github.com/projecteru2/core/utils"
)

const NodeResourceInfoKey = "/resource/cpumem/%s"

// GetNodeResourceInfo .
func (c *CPUMem) GetNodeResourceInfo(ctx context.Context, node string, workloadResourceMap map[string]*types.WorkloadResourceArgs, fix bool) (*types.NodeResourceInfo, []string, error) {
	resourceInfo, err := c.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		return nil, nil, err
	}

	diffs := []string{}

	totalResourceArgs := &types.WorkloadResourceArgs{}
	for _, args := range workloadResourceMap {
		totalResourceArgs.Add(args)
	}

	totalResourceArgs.CPURequest = utils.Round(totalResourceArgs.CPURequest)
	totalCPUUsage := utils.Round(resourceInfo.Usage.CPU)
	if totalResourceArgs.CPURequest != totalCPUUsage {
		diffs = append(diffs, fmt.Sprintf("node.CPUUsed != sum(workload.CPURequest): %.2f != %.2f", totalCPUUsage, totalResourceArgs.CPURequest))
	}

	for cpu := range totalResourceArgs.CPUMap {
		if totalResourceArgs.CPUMap[cpu] != resourceInfo.Usage.CPUMap[cpu] {
			diffs = append(diffs, fmt.Sprintf("node.CPUMap[%v] != sum(workload.CPUMap[%v]): %v != %v", cpu, cpu, resourceInfo.Usage.CPUMap[cpu], totalResourceArgs.CPUMap[cpu]))
		}
	}

	for numaNodeID := range totalResourceArgs.NUMAMemory {
		if totalResourceArgs.NUMAMemory[numaNodeID] != resourceInfo.Usage.NUMAMemory[numaNodeID] {
			diffs = append(diffs, fmt.Sprintf("node.NUMAMemory[%v] != sum(workload.NUMAMemory[%v]: %v != %v)", numaNodeID, numaNodeID, resourceInfo.Usage.NUMAMemory[numaNodeID], totalResourceArgs.NUMAMemory[numaNodeID]))
		}
	}

	if resourceInfo.Usage.Memory != totalResourceArgs.MemoryRequest {
		diffs = append(diffs, fmt.Sprintf("node.MemoryUsed != sum(workload.MemoryRequest): %d != %d", resourceInfo.Usage.Memory, totalResourceArgs.MemoryRequest))
	}

	if fix {
		resourceInfo.Usage = &types.NodeResourceArgs{
			CPUMap:     totalResourceArgs.CPUMap,
			Memory:     totalResourceArgs.MemoryRequest,
			NUMAMemory: totalResourceArgs.NUMAMemory,
		}
		if err = c.doSetNodeResourceInfo(ctx, node, resourceInfo); err != nil {
			logrus.Warnf("[GetNodeResourceInfo] failed to fix node resource, err: %v", err)
			diffs = append(diffs, "fix failed")
		}
	}

	return resourceInfo, diffs, nil
}

// SetNodeResourceInfo .
func (c *CPUMem) SetNodeResourceInfo(ctx context.Context, node string, resourceCapacity *types.NodeResourceArgs, resourceUsage *types.NodeResourceArgs) error {
	resourceInfo, err := c.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[SetNodeResourceInfo] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	resourceInfo.Capacity = resourceCapacity
	resourceInfo.Usage = resourceUsage

	return c.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

func (c *CPUMem) doGetNodeResourceInfo(ctx context.Context, node string) (*types.NodeResourceInfo, error) {
	resourceInfo := &types.NodeResourceInfo{}
	resp, err := c.store.GetOne(ctx, fmt.Sprintf(NodeResourceInfoKey, node))
	if err != nil {
		logrus.Errorf("[doGetNodeResourceInfo] failed to get node resource info, err: %v", err)
		return nil, err
	}
	if err = json.Unmarshal(resp.Value, resourceInfo); err != nil {
		logrus.Errorf("[doGetNodeResourceInfo] failed to unmarshal node resource info, err: %v", err)
		return nil, err
	}
	return resourceInfo, nil
}

func (c *CPUMem) doSetNodeResourceInfo(ctx context.Context, node string, resourceInfo *types.NodeResourceInfo) error {
	if err := resourceInfo.Validate(); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] invalid resource info %+v, err: %v", resourceInfo, err)
		return err
	}

	data, err := json.Marshal(resourceInfo)
	if err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to marshal resource info %+v, err: %v", resourceInfo, err)
		return err
	}

	if _, err = c.store.Put(ctx, fmt.Sprintf(NodeResourceInfoKey, node), string(data)); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to put resource info %+v, err: %v", resourceInfo, err)
		return err
	}
	return nil
}
