package schedule

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/projecteru2/core-plugins/volume/types"
)

func TestGetVolumePlans(t *testing.T) {
	var resourceInfo *types.NodeResourceInfo
	var request types.VolumeBindings
	var exists types.VolumeMap
	var maxDeployCount = 100
	var plans []types.VolumePlan
	var err error

	// normal case
	resourceInfo = &types.NodeResourceInfo{
		Capacity: &types.NodeResourceArgs{
			Volumes: map[string]int64{
				"/v1": 1000000000000,
				"/v2": 10000,
				"/v3": 10000,
			},
		},
		Usage: &types.NodeResourceArgs{
			Volumes: map[string]int64{
				"/v1": 0,
				"/v2": 0,
				"/v3": 1,
			},
		},
	}
	request, err = types.NewVolumeBindings([]string{
		"AUTO:/data0:rwm:100",
		"AUTO:/data1:rw:1",
		"AUTO:/data2:rw:0",
	})
	assert.Nil(t, err)
	plans = GetVolumePlans(resourceInfo, request, exists, maxDeployCount)
	fmt.Printf("total %v plans:\n", len(plans))
	for _, plan := range plans {
		body, _ := plan.MarshalJSON()
		fmt.Println(string(body))
	}
	fmt.Println("--------------------------")

	// multiple normal request
	resourceInfo = &types.NodeResourceInfo{
		Capacity: &types.NodeResourceArgs{
			Volumes: map[string]int64{
				"/v1": 1000000000000,
				"/v2": 5000,
				"/v3": 10000,
			},
		},
		Usage: &types.NodeResourceArgs{
			Volumes: map[string]int64{
				"/v1": 1,
				"/v2": 1,
				"/v3": 1,
			},
		},
	}
	request, err = types.NewVolumeBindings([]string{
		"AUTO:/data0:rw:100",
		"AUTO:/data1:rw:100",
		"AUTO:/data2:rw:100",
		"AUTO:/data3:rw:0",
	})
	assert.Nil(t, err)
	plans = GetVolumePlans(resourceInfo, request, exists, maxDeployCount)
	fmt.Printf("total %v plans:\n", len(plans))
	for _, plan := range plans {
		body, _ := plan.MarshalJSON()
		fmt.Println(string(body))
	}
	fmt.Println("--------------------------")

	resourceInfo = &types.NodeResourceInfo{
		Capacity: &types.NodeResourceArgs{
			Volumes: map[string]int64{
				"/v1": 10000,
				"/v2": 5000,
				"/v3": 10000,
			},
		},
		Usage: &types.NodeResourceArgs{
			Volumes: map[string]int64{
				"/v1": 0,
				"/v2": 0,
				"/v3": 1,
			},
		},
	}
	request, err = types.NewVolumeBindings([]string{
		"AUTO:/data0:rwm:100",
		"AUTO:/data1:rw:100",
		"AUTO:/data2:rw:0",
	})
	assert.Nil(t, err)
	plans = GetVolumePlans(resourceInfo, request, exists, maxDeployCount)
	fmt.Printf("total %v plans:\n", len(plans))
	for _, plan := range plans {
		body, _ := plan.MarshalJSON()
		fmt.Println(string(body))
	}
	fmt.Println("--------------------------")
}
