package models

import (
	"context"
	"errors"
	"testing"

	"github.com/docker/go-units"
	"github.com/stretchr/testify/assert"

	"github.com/projecteru2/core-plugins/cpumem/types"
	coretypes "github.com/projecteru2/core/types"
)

func TestGetNodeResourceInfo(t *testing.T) {
	ctx := context.Background()

	cpuMem := newTestCPUMem(t)
	nodes := generateNodes(t, cpuMem, 1, 2, 4*units.GiB, 100)
	node := nodes[0]

	// invalid node
	_, _, err := cpuMem.GetNodeResourceInfo(ctx, "xxx", nil, false)
	assert.True(t, errors.Is(err, coretypes.ErrBadCount))

	resourceInfo, diffs, err := cpuMem.GetNodeResourceInfo(ctx, node, nil, false)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(diffs))

	resourceInfo.Capacity.NUMA = types.NUMA{"0": "0", "1": "1"}
	resourceInfo.Capacity.NUMAMemory = types.NUMAMemory{"0": 2 * units.GiB, "1": 2 * units.GiB}

	assert.Nil(t, cpuMem.SetNodeResourceInfo(ctx, node, resourceInfo.Capacity, resourceInfo.Usage))

	resourceInfo, diffs, err = cpuMem.GetNodeResourceInfo(ctx, node, map[string]*types.WorkloadResourceArgs{
		"x-workload": {
			CPURequest:    2,
			CPUMap:        types.CPUMap{"0": 100, "1": 100},
			MemoryRequest: 2 * units.GiB,
			NUMAMemory:    types.NUMAMemory{"0": units.GiB, "1": units.GiB},
		},
	}, true)
	assert.Nil(t, err)
	assert.Equal(t, 6, len(diffs))
	assert.Equal(t, resourceInfo.Usage, &types.NodeResourceArgs{
		CPU:        2,
		CPUMap:     types.CPUMap{"0": 100, "1": 100},
		Memory:     2 * units.GiB,
		NUMAMemory: types.NUMAMemory{"0": units.GiB, "1": units.GiB},
	})
}

func TestSetNodeResourceInfo(t *testing.T) {
	ctx := context.Background()

	cpuMem := newTestCPUMem(t)
	nodes := generateNodes(t, cpuMem, 1, 2, 4*units.GiB, 100)
	node := nodes[0]

	resourceInfo, _, err := cpuMem.GetNodeResourceInfo(ctx, node, nil, false)
	assert.Nil(t, err)

	assert.True(t, errors.Is(cpuMem.SetNodeResourceInfo(ctx, "xxx", resourceInfo.Capacity, resourceInfo.Usage), coretypes.ErrBadCount))
}