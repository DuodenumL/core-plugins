package types

import (
	"encoding/json"
	"testing"

	"github.com/docker/go-units"
	"github.com/stretchr/testify/assert"
)

func TestWorkloadResourceOptsValidate(t *testing.T) {
	var err error
	var resourceOpts *WorkloadResourceOpts
	// Mem request below zero shall fail
	resourceOpts = &WorkloadResourceOpts{
		MemRequest: -1,
		MemLimit:   -1,
	}
	err = resourceOpts.Validate()
	assert.NotNil(t, err)

	// Mem and cpu request equal to zero will not fail
	resourceOpts = &WorkloadResourceOpts{
		MemRequest: 0,
		MemLimit:   1,
		CPURequest: 0,
		CPULimit:   1,
	}
	err = resourceOpts.Validate()
	assert.Nil(t, err)

	// Request more then limited will not fail
	resourceOpts = &WorkloadResourceOpts{
		MemRequest: 2,
		MemLimit:   1,
		CPURequest: 2,
		CPULimit:   1,
	}
	err = resourceOpts.Validate()
	assert.Nil(t, err)

	// Request below zero will fail
	resourceOpts = &WorkloadResourceOpts{
		CPURequest: -0.5,
		CPULimit:   -1,
	}
	err = resourceOpts.Validate()
	assert.NotNil(t, err)

	// Request unlimited cpu but with cpu bind will fail
	resourceOpts = &WorkloadResourceOpts{
		CPURequest: 0,
		CPUBind:    true,
	}
	err = resourceOpts.Validate()
	assert.NotNil(t, err)

	// Request cpu bind and limit < request
	resourceOpts = &WorkloadResourceOpts{
		CPUBind: true,
		CPURequest: 2.2,
		CPULimit: 3.3,
	}
	err = resourceOpts.Validate()
	assert.Nil(t, err)
	assert.Equal(t, resourceOpts.CPURequest, 3.3)
	assert.Equal(t, resourceOpts.CPULimit, 3.3)
}

func rawParamsToWorkloadResourceOpts(r RawParams, w *WorkloadResourceOpts) error {
	var body []byte
	var err error
	if body, err = json.Marshal(r); err != nil {
		return err
	}
	if err = w.ParseFromString(string(body)); err != nil {
		return err
	}
	return w.Validate()
}

func TestWorkloadResourceOptsParseFromString(t *testing.T) {
	resourceOpts := &WorkloadResourceOpts{}

	// invalid json
	invalidJsonStr := "xxx"
	assert.NotNil(t, resourceOpts.ParseFromString(invalidJsonStr))

	// normal case
	r := RawParams{
		"cpu-bind":       nil,
		"cpu-request":    1.1,
		"cpu-limit":      2.2,
		"cpu":            3.3,
		"memory-request": "1GB",
		"memory-limit":   "2GB",
		"memory":         "3GB",
	}
	assert.Nil(t, rawParamsToWorkloadResourceOpts(r, resourceOpts))
	assert.Equal(t, resourceOpts, &WorkloadResourceOpts{
		CPUBind:    true,
		CPURequest: 3.3,
		CPULimit:   3.3,
		MemRequest: 3 * units.GiB,
		MemLimit:   3 * units.GiB,
	})

	// no cpu shortcut
	delete(r, "cpu")
	delete(r, "cpu-bind")
	assert.Nil(t, rawParamsToWorkloadResourceOpts(r, resourceOpts))
	assert.Equal(t, resourceOpts.CPURequest, 1.1)
	assert.Equal(t, resourceOpts.CPULimit, 2.2)

	// no memory shortcut
	delete(r, "memory")
	assert.Nil(t, rawParamsToWorkloadResourceOpts(r, resourceOpts))
	assert.EqualValues(t, resourceOpts.MemRequest, units.GiB)
	assert.EqualValues(t, resourceOpts.MemLimit, 2*units.GiB)
}

func rawParamsToNodeResourceOpts(r RawParams, n *NodeResourceOpts) error {
	var body []byte
	var err error
	if body, err = json.Marshal(r); err != nil {
		return err
	}
	return n.ParseFromString(string(body))
}

func TestNodeResourceOptsParseFromString(t *testing.T) {
	resourceOpts := &NodeResourceOpts{}

	// invalid json
	invalidJsonStr := "xxx"
	assert.NotNil(t, resourceOpts.ParseFromString(invalidJsonStr))

	// invalid cpu id / pieces
	r := RawParams{
		"cpu": "0:100,1:100,xxx:100",
	}
	assert.NotNil(t, rawParamsToNodeResourceOpts(r, resourceOpts))
	r = RawParams{
		"cpu": "0:100,1:100,2:100x",
	}
	assert.NotNil(t, rawParamsToNodeResourceOpts(r, resourceOpts))

	// invalid memory
	r = RawParams{
		"cpu": "4",
		"memory": "1kg",
	}
	assert.NotNil(t, rawParamsToNodeResourceOpts(r, resourceOpts))

	// invalid numa memory
	r = RawParams{
		"cpu": "4",
		"memory": "1G",
		"numa-cpu": []string{"0,1","2,3"},
		"numa-memory": []string{"1kg", "1kg"},
	}
	assert.NotNil(t, rawParamsToNodeResourceOpts(r, resourceOpts))

	// normal case
	r = RawParams{
		"cpu": "4",
		"memory": "2G",
		"numa-cpu": []string{"0,1","2,3"},
		"numa-memory": []string{"1G", "1G"},
	}
	assert.Nil(t, rawParamsToNodeResourceOpts(r, resourceOpts))
	assert.Equal(t, resourceOpts, &NodeResourceOpts{
		CPUMap: CPUMap{"0": 100, "1": 100, "2": 100, "3": 100},
		Memory:     2 * units.GiB,
		NUMA:       NUMA{"0": "0", "1": "0", "2": "1", "3": "1"},
		NUMAMemory: NUMAMemory{"0": units.GiB, "1": units.GiB},
	})
}