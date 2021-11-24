package models

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/docker/go-units"
	"github.com/stretchr/testify/assert"

	"github.com/projecteru2/core-plugins/cpumem/types"
	"github.com/projecteru2/core/store/etcdv3/meta"
	coretypes "github.com/projecteru2/core/types"
)

func generateNodeResourceInfos(t *testing.T, nums int, cores int, memory int64, shares int) []*types.NodeResourceInfo {
	infos := []*types.NodeResourceInfo{}
	for i := 0; i < nums; i++ {
		cpuMap := types.CPUMap{}
		for c := 0; c < cores; c++ {
			cpuMap[strconv.Itoa(c)] = shares
		}

		info := &types.NodeResourceInfo{
			Capacity: &types.NodeResourceArgs{
				CPU:    float64(cores),
				CPUMap: cpuMap,
				Memory: memory,
			},
			Usage: nil,
		}
		assert.Nil(t, info.Validate())

		infos = append(infos, info)
	}
	return infos
}

func generateNodes(t *testing.T, cpuMem *CPUMem, nums int, cores int, memory int64, shares int) []string {
	nodes := []string{}
	infos := generateNodeResourceInfos(t, nums, cores, memory, shares)

	for i, info := range infos {
		nodeName := fmt.Sprintf("node%d", i)
		assert.Nil(t, cpuMem.doSetNodeResourceInfo(context.Background(), nodeName, info))
		nodes = append(nodes, nodeName)
	}

	return nodes
}

func generateComplexNodes(t *testing.T, cpuMem *CPUMem) []string {
	infos := []*types.NodeResourceInfo{
		{
			Capacity: &types.NodeResourceArgs{
				CPU: 4,
				CPUMap: types.CPUMap{
					"0": 100,
					"1": 100,
					"2": 100,
					"3": 100,
				},
				Memory: 12 * units.GiB,
			},
		},
		{
			Capacity: &types.NodeResourceArgs{
				CPU: 14,
				CPUMap: types.CPUMap{
					"0":  100,
					"1":  100,
					"10": 100,
					"11": 100,
					"12": 100,
					"13": 100,
					"2":  100,
					"3":  100,
					"4":  100,
					"5":  100,
					"6":  100,
					"7":  100,
					"8":  100,
					"9":  100,
				},
				Memory: 12 * units.GiB,
			},
		},
		{
			Capacity: &types.NodeResourceArgs{
				CPU: 12,
				CPUMap: types.CPUMap{
					"0":  100,
					"1":  100,
					"10": 100,
					"11": 100,
					"2":  100,
					"3":  100,
					"4":  100,
					"5":  100,
					"6":  100,
					"7":  100,
					"8":  100,
					"9":  100,
				},
				Memory: 12 * units.GiB,
			},
		},
		{
			Capacity: &types.NodeResourceArgs{
				CPU: 18,
				CPUMap: types.CPUMap{
					"0":  100,
					"1":  100,
					"10": 100,
					"11": 100,
					"12": 100,
					"13": 100,
					"14": 100,
					"15": 100,
					"16": 100,
					"17": 100,
					"2":  100,
					"3":  100,
					"4":  100,
					"5":  100,
					"6":  100,
					"7":  100,
					"8":  100,
					"9":  100,
				},
				Memory: 12 * units.GiB,
			},
		},
		{
			Capacity: &types.NodeResourceArgs{
				CPU: 8,
				CPUMap: types.CPUMap{
					"0": 100,
					"1": 100,
					"2": 100,
					"3": 100,
					"4": 100,
					"5": 100,
					"6": 100,
					"7": 100,
				},
				Memory: 12 * units.GiB,
			},
		},
	}
	nodes := []string{}
	for i, info := range infos {
		nodeName := fmt.Sprintf("node%d", i)
		assert.Nil(t, cpuMem.doSetNodeResourceInfo(context.Background(), nodeName, info))
		nodes = append(nodes, nodeName)
	}
	return nodes
}

func newTestCPUMem(t *testing.T) *CPUMem {
	config := &types.Config{
		Scheduler: types.SchedConfig{
			MaxShare:  -1,
			ShareBase: 100,
		},
	}
	cpuMem := &CPUMem{
		config: config,
	}
	store, err := meta.NewETCD(coretypes.EtcdConfig{Prefix: "/cpumem"}, t)
	assert.Nil(t, err)
	cpuMem.store = store
	return cpuMem
}

func TestGetNodesCapacityWithCPUBinding(t *testing.T) {
	ctx := context.Background()

	cpuMem := newTestCPUMem(t)
	nodes := generateNodes(t, cpuMem, 2, 2, 4*units.GiB, 100)

	_, total, err := cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 0.5,
		MemRequest: 1,
	})
	assert.Nil(t, err)
	assert.True(t, total >= 1)

	_, total, err = cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 2,
		MemRequest: 1,
	})
	assert.Nil(t, err)
	assert.True(t, total < 3)

	_, total, err = cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 3,
		MemRequest: 1,
	})
	assert.Nil(t, err)
	assert.True(t, total < 2)

	_, total, err = cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 1,
		MemRequest: 1,
	})
	assert.Nil(t, err)
	assert.True(t, total < 5)
}

func TestComplexNodes(t *testing.T) {
	ctx := context.Background()

	cpuMem := newTestCPUMem(t)
	nodes := generateComplexNodes(t, cpuMem)
	_, total, err := cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 1.7,
		MemRequest: 1,
	})
	assert.Nil(t, err)
	assert.Equal(t, 28, total)
}

func TestCPUNodesWithMemoryLimit(t *testing.T) {
	ctx := context.Background()

	cpuMem := newTestCPUMem(t)
	nodes := generateNodes(t, cpuMem, 2, 2, 1024, 100)
	_, total, err := cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 0.1,
		MemRequest: 1024,
	})
	assert.Nil(t, err)
	assert.Equal(t, total, 2)

	_, total, err = cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 0.1,
		MemRequest: 1025,
	})
	assert.Nil(t, err)
	assert.Equal(t, total, 0)
}

func TestCPUNodesWithMaxShareLimit(t *testing.T) {
	ctx := context.Background()

	cpuMem := newTestCPUMem(t)
	cpuMem.config.Scheduler.MaxShare = 2

	nodes := generateNodes(t, cpuMem, 1, 6, 12*units.GiB, 100)
	_, total, err := cpuMem.GetNodesCapacity(ctx, nodes, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 1.7,
		MemRequest: 1,
	})
	assert.Nil(t, err)
	assert.Equal(t, total, 2)

	nodeResourceInfo := &types.NodeResourceInfo{Capacity: &types.NodeResourceArgs{
		CPU:    4,
		CPUMap: types.CPUMap{"0": 0, "1": 0, "2": 100, "3": 100},
		Memory: 12 * units.GiB,
	}}
	assert.Nil(t, nodeResourceInfo.Validate())
	_, _, err = cpuMem.doAllocByCPU(nodeResourceInfo, 1, &types.WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 1.2,
		MemRequest: 1,
	})
	assert.Nil(t, err)
}

func BenchmarkGetNodesCapacity(b *testing.B) {
	b.StopTimer()
	t := &testing.T{}
	cpuMem := newTestCPUMem(t)
	nodes := generateNodes(t, cpuMem, 10000, 24, 128*units.GiB, 100)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := cpuMem.GetNodesCapacity(context.Background(), nodes, &types.WorkloadResourceOpts{
			CPUBind:    true,
			CPURequest: 1.3,
			MemRequest: 1,
		})
		assert.Nil(b, err)
	}
}
