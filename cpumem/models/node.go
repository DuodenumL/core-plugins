package models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/projecteru2/core-plugins/cpumem/types"
	coretypes "github.com/projecteru2/core/types"
	"github.com/sirupsen/logrus"
)

const (
	NodeResourceInfoKey = "/resource/cpumem/%s"
	Incr                = true
	Decr                = false
)

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

	totalResourceArgs.CPURequest = coretypes.Round(totalResourceArgs.CPURequest)
	totalCPUUsage := coretypes.Round(resourceInfo.Usage.CPU)
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
func (c *CPUMem) SetNodeResourceInfo(ctx context.Context, node string, resourceOpts *types.NodeResourceOpts) (*types.NodeResourceInfo, error) {
	resourceInfo := &types.NodeResourceInfo{
		Capacity: &types.NodeResourceArgs{
			CPU:        float64(len(resourceOpts.CPUMap)),
			CPUMap:     resourceOpts.CPUMap,
			Memory:     resourceOpts.Memory,
			NUMAMemory: resourceOpts.NUMAMemory,
		},
		Usage: &types.NodeResourceArgs{
			CPU:        float64(len(resourceOpts.CPUMap)),
			CPUMap:     types.CPUMap{},
			Memory:     0,
			NUMAMemory: types.NUMAMemory{},
		},
		NUMA: resourceOpts.NUMA,
	}

	// if NUMA is set but NUMAMemory is not set
	// then divide memory equally according to the number of numa nodes
	if resourceOpts.NUMA != nil && resourceOpts.NUMAMemory == nil {
		averageMemory := resourceOpts.Memory / int64(len(resourceOpts.NUMA))
		resourceInfo.Capacity.NUMAMemory = types.NUMAMemory{}
		for _, numaNodeID := range resourceOpts.NUMA {
			resourceInfo.Capacity.NUMAMemory[numaNodeID] = averageMemory
		}
	}

	for cpu := range resourceOpts.CPUMap {
		resourceInfo.Usage.CPUMap[cpu] = 0
	}

	for numaNodeID := range resourceOpts.NUMA {
		resourceInfo.Usage.NUMAMemory[numaNodeID] = 0
	}

	return resourceInfo, c.doSetNodeResourceInfo(ctx, node, resourceInfo)
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

// UpdateNodeResourceUsage .
func (c *CPUMem) UpdateNodeResourceUsage(ctx context.Context, node string, resourceArgs *types.NodeResourceArgs, direction bool) error {
	resourceInfo, err := c.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[UpdateNodeResourceUsage] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	if direction == Incr {
		resourceInfo.Usage.Add(resourceArgs)
	} else {
		resourceInfo.Usage.Sub(resourceArgs)
	}



	return c.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

// UpdateNodeResourceCapacity .
func (c *CPUMem) UpdateNodeResourceCapacity(ctx context.Context, node string, resourceOpts *types.NodeResourceOpts, direction bool) error {
	resourceInfo, err := c.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[UpdateNodeResourceCapacity] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	resourceArgs := &types.NodeResourceArgs{
		CPUMap:     resourceOpts.CPUMap,
		Memory:     resourceOpts.Memory,
		NUMAMemory: resourceOpts.NUMAMemory,
	}

	if len(resourceOpts.NUMA) > 0 {
		resourceInfo.NUMA = resourceOpts.NUMA
	}

	if direction == Incr {
		resourceInfo.Capacity.Add(resourceArgs)
	} else {
		resourceInfo.Capacity.Sub(resourceArgs)
	}

	// add new cpu
	for cpu := range resourceInfo.Capacity.CPUMap {
		_, ok := resourceInfo.Usage.CPUMap[cpu]
		if !ok {
			resourceInfo.Usage.CPUMap[cpu] = 0
			continue
		}
	}

	// delete cpus with no pieces
	resourceInfo.RemoveEmptyCores()

	return c.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

// AddNode .
func (c *CPUMem) AddNode(ctx context.Context, node string, resourceOpts *types.NodeResourceOpts) (*types.NodeResourceInfo, error) {
	if _, err := c.doGetNodeResourceInfo(ctx, node); err != nil {
		if errors.Is(err, coretypes.ErrBadCount) {
			return nil, types.ErrNodeExists
		}
		logrus.Errorf("[AddNode] failed to get resource info of node %v, err: %v", node, err)
		return nil, err
	}

	return c.SetNodeResourceInfo(ctx, node, resourceOpts)
}

// RemoveNode .
func (c *CPUMem) RemoveNode(ctx context.Context, node string) error {
	if _, err := c.store.Delete(ctx, fmt.Sprintf(NodeResourceInfoKey, node)); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to delete node %v, err: %v", node, err)
		return err
	}
	return nil
}