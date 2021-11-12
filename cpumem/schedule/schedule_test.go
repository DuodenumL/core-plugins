package schedule

import (
	"github.com/projecteru2/core-plugins/cpumem/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func assertCPUPlansEqual(t *testing.T, m1 []types.CPUMap, m2 []types.CPUMap) {
	// TODO
	assert.Equal(t, len(m1), len(m2))
	for i := range m1 {
		assert.Equal(t, m1[i].TotalPieces(), m2[i].TotalPieces())
	}
}

func TestGetCPUPlans(t *testing.T) {
	var cpuMap types.CPUMap
	var cpuRequest float64
	var cpuPlans []types.CPUMap
	var expectedCPUPlans []types.CPUMap
	shareBase := 100
	maxFragmentCores := -1
	memRequest := int64(0)
	availableMemory := int64(10000000)

	// normal request
	cpuMap = types.CPUMap{"0": 100, "1": 100, "2": 100, "3": 100}
	cpuRequest = 1.2
	cpuPlans = doGetCPUPlans(nil, cpuMap, availableMemory, shareBase, maxFragmentCores, cpuRequest, memRequest)
	expectedCPUPlans = []types.CPUMap{
		{
			"0": 20,
			"1": 100,
		},
		{
			"0": 20,
			"2": 100,
		},
		{
			"0": 20,
			"3": 100,
		},
	}
	assertCPUPlansEqual(t, expectedCPUPlans, cpuPlans)

	// only two full cores
	cpuMap = types.CPUMap{"0": 90, "1": 90, "2": 100, "3": 100}
	cpuRequest = 1.2
	cpuPlans = doGetCPUPlans(nil, cpuMap, availableMemory, shareBase, maxFragmentCores, cpuRequest, memRequest)
	expectedCPUPlans = []types.CPUMap{
		{
			"0": 20,
			"2": 100,
		},
		{
			"0": 20,
			"3": 100,
		},
	}
	assertCPUPlansEqual(t, expectedCPUPlans, cpuPlans)

	// oversell
	cpuMap = types.CPUMap{"0": 200, "1": 200, "2": 200, "3": 200}
	cpuRequest = 2.2
	cpuPlans = doGetCPUPlans(nil, cpuMap, availableMemory, shareBase, maxFragmentCores, cpuRequest, memRequest)
	expectedCPUPlans = []types.CPUMap{
		{
			"1": 100,
			"2": 100,
			"0": 20,
		},

	}
	assertCPUPlansEqual(t, expectedCPUPlans, cpuPlans)

	// normal case
	cpuMap = types.CPUMap{"0": 200, "1": 200, "2": 200, "3": 200}
	cpuRequest = 1.8
	cpuPlans = doGetCPUPlans(nil, cpuMap, availableMemory, shareBase, maxFragmentCores, cpuRequest, memRequest)
	expectedCPUPlans = []types.CPUMap{
		{
			"0": 80,
			"2": 100,
		},
		{
			"0": 80,
			"3": 100,
		},
	}
	assertCPUPlansEqual(t, expectedCPUPlans, cpuPlans)

	// normal case
	cpuMap = types.CPUMap{"0": 100, "1": 100, "2": 100, "3": 100}
	cpuRequest = 1.8
	cpuPlans = doGetCPUPlans(nil, cpuMap, availableMemory, shareBase, maxFragmentCores, cpuRequest, memRequest)
	expectedCPUPlans = []types.CPUMap{
		{
			"0": 80,
			"2": 100,
		},
		{
			"1": 80,
			"3": 100,
		},
	}
	assertCPUPlansEqual(t, expectedCPUPlans, cpuPlans)

	// test 2
	cpuMap = types.CPUMap{"0": 0, "1": 0, "2": 100, "3": 100}
	cpuRequest = 1.2
	cpuPlans = doGetCPUPlans(nil, cpuMap, availableMemory, shareBase, 2, cpuRequest, memRequest)
	expectedCPUPlans = []types.CPUMap{
		{
			"0": 80,
			"2": 100,
		},
		{
			"0": 80,
			"3": 100,
		},
	}
	assert.True(t, len(cpuPlans) > 0)
}
