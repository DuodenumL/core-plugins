package types

// CPUMap .
type CPUMap map[string]int

// TotalPieces .
func (c CPUMap) TotalPieces() int {
	res := 0
	for _, pieces := range c {
		res += pieces
	}
	return res
}

// DeepCopy .
func (c CPUMap) DeepCopy() CPUMap {
	res := CPUMap{}
	for cpu, pieces := range c {
		res[cpu] = pieces
	}
	return res
}

// Sub .
func (c CPUMap) Sub(c1 CPUMap) {
	for cpu, pieces := range c1 {
		c[cpu] -= pieces
	}
}

// Add .
func (c CPUMap) Add(c1 CPUMap) {
	for cpu, pieces := range c1 {
		c[cpu] += pieces
	}
}

// NUMA map[cpuID]nodeID
type NUMA map[string]string

// NUMAMemory .
type NUMAMemory map[string]int64

// WorkloadResourceArgs .
type WorkloadResourceArgs struct {
	CPURequest    float64    `json:"cpu_request"`
	CPULimit      float64    `json:"cpu_limit"`
	MemoryRequest int64      `json:"memory_request"`
	MemoryLimit   int64      `json:"memory_limit"`
	CPUMap        CPUMap     `json:"cpu_map"`
	NUMAMemory    NUMAMemory `json:"numa_memory"`
	NUMANode      string     `json:"numa_node"`
}

func (r *WorkloadResourceArgs) Add(r1 *WorkloadResourceArgs) {
	r.CPURequest += r1.CPURequest
	r.MemoryRequest += r1.MemoryRequest
	if r.CPUMap == nil {
		r.CPUMap = r1.CPUMap
	} else {
		for cpu := range r1.CPUMap {
			r.CPUMap[cpu] += r1.CPUMap[cpu]
		}
	}
	if r.NUMAMemory == nil {
		r.NUMAMemory = r1.NUMAMemory
	} else {
		for nodeID := range r1.NUMAMemory {
			r.NUMAMemory[nodeID] += r1.NUMAMemory[nodeID]
		}
	}
}

// NodeResourceArgs .
type NodeResourceArgs struct {
	CPU        float64    `json:"cpu"`
	CPUMap     CPUMap     `json:"cpu_map"`
	Memory     int64      `json:"memory"`
	NUMAMemory NUMAMemory `json:"numa_memory"`
}

// DeepCopy .
func (r *NodeResourceArgs) DeepCopy() *NodeResourceArgs {
	res := &NodeResourceArgs{
		CPU:        r.CPU,
		CPUMap:     CPUMap{},
		Memory:     r.Memory,
		NUMAMemory: NUMAMemory{},
	}

	for cpu := range r.CPUMap {
		res.CPUMap[cpu] = r.CPUMap[cpu]
	}
	for numaNodeID := range r.NUMAMemory {
		res.NUMAMemory[numaNodeID] = r.NUMAMemory[numaNodeID]
	}
	return res
}

// Add .
func (r *NodeResourceArgs) Add(r1 *NodeResourceArgs) {
	r.CPU += r1.CPU
	r.CPUMap.Add(r1.CPUMap)
	r.Memory += r1.Memory

	for numaNodeID := range r1.NUMAMemory {
		r.NUMAMemory[numaNodeID] += r1.NUMAMemory[numaNodeID]
	}
}

// Sub .
func (r *NodeResourceArgs) Sub(r1 *NodeResourceArgs) {
	r.CPU -= r1.CPU
	r.CPUMap.Sub(r1.CPUMap)
	r.Memory -= r1.Memory

	for numaNodeID := range r1.NUMAMemory {
		r.NUMAMemory[numaNodeID] -= r1.NUMAMemory[numaNodeID]
	}
}

// NodeResourceInfo .
type NodeResourceInfo struct {
	Capacity *NodeResourceArgs
	Usage    *NodeResourceArgs
	NUMA     NUMA
}

// RemoveEmptyCores .
func (n *NodeResourceInfo) RemoveEmptyCores() {
	keysToDelete := []string{}
	for cpu := range n.Capacity.CPUMap {
		if n.Capacity.CPUMap[cpu] == 0 && n.Usage.CPUMap[cpu] == 0 {
			keysToDelete = append(keysToDelete, cpu)
		}
	}

	for _, cpu := range keysToDelete {
		delete(n.Capacity.CPUMap, cpu)
		delete(n.Usage.CPUMap, cpu)
	}
}

func (n *NodeResourceInfo) Validate() error {
	if n.Capacity == nil || n.Capacity.CPUMap == nil {
		return ErrInvalidCapacity
	}
	if n.Usage == nil || n.Usage.CPUMap == nil {
		return ErrInvalidUsage
	}
	if len(n.Capacity.CPUMap) == 0 || len(n.Capacity.CPUMap) != len(n.Usage.CPUMap) {
		return ErrInvalidCPUMap
	}

	for cpu, totalPieces := range n.Capacity.CPUMap {
		if totalPieces < 0 {
			return ErrInvalidCPUMap
		}
		if piecesUsed, ok := n.Usage.CPUMap[cpu]; !ok || piecesUsed < 0 || piecesUsed > totalPieces {
			return ErrInvalidCPUMap
		}
	}

	if n.NUMA != nil {
		if n.Capacity.NUMAMemory == nil || n.Usage.NUMAMemory == nil {
			return ErrInvalidNUMAMemory
		}
		for cpu := range n.Capacity.CPUMap {
			if _, ok := n.NUMA[cpu]; !ok {
				return ErrInvalidNUMA
			}
		}

		for numaNodeID, totalMemory := range n.Capacity.NUMAMemory {
			if totalMemory < 0 {
				return ErrInvalidNUMAMemory
			}
			if memoryUsed, ok := n.Usage.NUMAMemory[numaNodeID]; !ok || memoryUsed < 0 || memoryUsed > totalMemory {
				return ErrInvalidNUMAMemory
			}
		}
	}

	return nil
}

func (n *NodeResourceInfo) GetAvailableResource() *NodeResourceArgs {
	availableResourceArgs := n.Capacity.DeepCopy()
	availableResourceArgs.Sub(n.Usage)

	return availableResourceArgs
}

// WorkloadResourceOpts includes all possible fields passed by eru-core for editing workload
type WorkloadResourceOpts struct {
	CPUBind    bool    `json:"cpu_bind"`
	CPURequest float64 `json:"cpu_request"`
	CPULimit   float64 `json:"cpu_limit"`
	MemRequest int64   `json:"mem_request"`
	MemLimit   int64   `json:"mem_limit"`
}

// NodeResourceOpts includes all possible fields passed by eru-core for editing node
type NodeResourceOpts struct {
	CPUMap     CPUMap     `json:"cpu_map"`
	Memory     int64      `json:"memory"`
	NUMA       NUMA       `json:"numa"`
	NUMAMemory NUMAMemory `json:"numa_memory"`
}

// NodeCapacityInfo .
type NodeCapacityInfo struct {
	Node     string  `json:"node"`
	Capacity int     `json:"capacity"`
	Usage    float64 `json:"usage"`
	Rate     float64 `json:"rate"`
	Weight   int     `json:"weight"`
}

// EngineArgs .
type EngineArgs struct {
	CPU      float64 `json:"cpu"`
	CPUMap   CPUMap  `json:"cpu_map"`
	NUMANode string  `json:"numa_node"`
	Memory   int64   `json:"memory"`
}
