package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"github.com/projecteru2/core-plugins/cpumem/utils"
	coreutils "github.com/projecteru2/core/utils"
)

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

// Add .
func (n NUMAMemory) Add(n1 NUMAMemory) {
	for numaNodeID, memory := range n1 {
		n[numaNodeID] += memory
	}
}

// Sub .
func (n NUMAMemory) Sub(n1 NUMAMemory) {
	for numaNodeID, memory := range n1 {
		n[numaNodeID] -= memory
	}
}

// RawParams .
type RawParams map[string]interface{}

// IsSet .
func (r RawParams) IsSet(key string) bool {
	_, ok := r[key]
	return ok
}

// Float64 .
func (r RawParams) Float64(key string) float64 {
	res, _ := strconv.ParseFloat(fmt.Sprintf("%v", r[key]), 64)
	return res
}

// Int64 .
func (r RawParams) Int64(key string) int64 {
	res, _ := strconv.ParseInt(fmt.Sprintf("%v", r[key]), 10, 64)
	return res
}

// String .
func (r RawParams) String(key string) string {
	if !r.IsSet(key) {
		return ""
	}
	if str, ok := r[key].(string); ok {
		return str
	}
	return ""
}

// StringSlice .
func (r RawParams) StringSlice(key string) []string {
	if !r.IsSet(key) {
		return nil
	}
	res := []string{}
	if s, ok := r[key].([]interface{}); ok {
		for _, v := range s {
			if str, ok := v.(string); ok {
				res = append(res, str)
			} else {
				return nil
			}
		}
	}
	return res
}

// Bool .
func (r RawParams) Bool(key string) bool {
	return r.IsSet(key)
}

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

// DeepCopy .
func (r *WorkloadResourceArgs) DeepCopy() *WorkloadResourceArgs {
	res := &WorkloadResourceArgs{
		CPURequest:    r.CPURequest,
		CPULimit:      r.CPULimit,
		MemoryRequest: r.MemoryRequest,
		MemoryLimit:   r.MemoryLimit,
		CPUMap:        CPUMap{},
		NUMAMemory:    NUMAMemory{},
		NUMANode:      r.NUMANode,
	}

	for cpu, pieces := range r.CPUMap {
		res.CPUMap[cpu] = pieces
	}
	for cpuID, numaNodeID := range res.NUMAMemory {
		res.NUMAMemory[cpuID] = numaNodeID
	}

	return res
}

// Add .
func (r *WorkloadResourceArgs) Add(r1 *WorkloadResourceArgs) {
	r.CPURequest = coreutils.Round(r.CPURequest + r1.CPURequest)
	r.MemoryRequest += r1.MemoryRequest
	if len(r.CPUMap) == 0 {
		r.CPUMap = r1.CPUMap
	} else {
		for cpu := range r1.CPUMap {
			r.CPUMap[cpu] += r1.CPUMap[cpu]
		}
	}
	if len(r.NUMAMemory) == 0 {
		r.NUMAMemory = r1.NUMAMemory
	} else {
		r.NUMAMemory.Add(r1.NUMAMemory)
	}
}

// Sub .
func (r *WorkloadResourceArgs) Sub(r1 *WorkloadResourceArgs) {
	r.CPURequest = coreutils.Round(r.CPURequest - r1.CPURequest)
	r.MemoryRequest -= r1.MemoryRequest
	if len(r.CPUMap) == 0 {
		r.CPUMap = CPUMap{}
	}
	r.CPUMap.Sub(r1.CPUMap)
	if r.NUMAMemory == nil {
		r.NUMAMemory = NUMAMemory{}
	}
	r.NUMAMemory.Sub(r1.NUMAMemory)
}

// NodeResourceArgs .
type NodeResourceArgs struct {
	CPU        float64    `json:"cpu"`
	CPUMap     CPUMap     `json:"cpu_map"`
	Memory     int64      `json:"memory"`
	NUMAMemory NUMAMemory `json:"numa_memory"`
	NUMA       NUMA       `json:"numa,omitempty"`
}

// DeepCopy .
func (r *NodeResourceArgs) DeepCopy() *NodeResourceArgs {
	res := &NodeResourceArgs{
		CPU:        r.CPU,
		CPUMap:     CPUMap{},
		Memory:     r.Memory,
		NUMAMemory: NUMAMemory{},
		NUMA:       NUMA{},
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
	r.CPU = coreutils.Round(r.CPU + r1.CPU)
	r.CPUMap.Add(r1.CPUMap)
	r.Memory += r1.Memory

	for numaNodeID := range r1.NUMAMemory {
		r.NUMAMemory[numaNodeID] += r1.NUMAMemory[numaNodeID]
	}
}

// Sub .
func (r *NodeResourceArgs) Sub(r1 *NodeResourceArgs) {
	r.CPU = coreutils.Round(r.CPU - r1.CPU)
	r.CPUMap.Sub(r1.CPUMap)
	r.Memory -= r1.Memory

	for numaNodeID := range r1.NUMAMemory {
		r.NUMAMemory[numaNodeID] -= r1.NUMAMemory[numaNodeID]
	}
}

// NodeResourceInfo .
type NodeResourceInfo struct {
	Capacity *NodeResourceArgs `json:"capacity"`
	Usage    *NodeResourceArgs `json:"usage"`
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

	n.Capacity.CPU = float64(len(n.Capacity.CPUMap))
}

func (n *NodeResourceInfo) Validate() error {
	if n.Capacity == nil || len(n.Capacity.CPUMap) == 0 {
		return ErrInvalidCapacity
	}
	if n.Usage == nil {
		n.Usage = &NodeResourceArgs{
			CPU:        0,
			CPUMap:     CPUMap{},
			Memory:     0,
			NUMAMemory: NUMAMemory{},
		}
		for cpuID := range n.Capacity.CPUMap {
			n.Usage.CPUMap[cpuID] = 0
		}
		for numaNodeID := range n.Capacity.NUMAMemory {
			n.Usage.NUMAMemory[numaNodeID] = 0
		}
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

	if len(n.Capacity.NUMA) > 0 {
		for cpu := range n.Capacity.CPUMap {
			if numaNodeID, ok := n.Capacity.NUMA[cpu]; !ok {
				return ErrInvalidNUMA
			} else if _, ok = n.Capacity.NUMAMemory[numaNodeID]; !ok {
				return ErrInvalidNUMAMemory
			}
		}

		for numaNodeID, totalMemory := range n.Capacity.NUMAMemory {
			if totalMemory < 0 {
				return ErrInvalidNUMAMemory
			}
			if memoryUsed := n.Usage.NUMAMemory[numaNodeID]; memoryUsed < 0 || memoryUsed > totalMemory {
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
	CPUBind     bool    `json:"cpu_bind"`
	KeepCPUBind bool    `json:"keep_cpu_bind"`
	CPURequest  float64 `json:"cpu_request"`
	CPULimit    float64 `json:"cpu_limit"`
	MemRequest  int64   `json:"mem_request"`
	MemLimit    int64   `json:"mem_limit"`
}

// Validate .
func (w *WorkloadResourceOpts) Validate() error {
	if w.CPURequest == 0 && w.CPULimit > 0 {
		w.CPURequest = w.CPULimit
	}
	if w.MemLimit < 0 || w.MemRequest < 0 {
		return errors.Wrap(ErrInvalidMemory, "limit or request less than 0")
	}
	if w.CPURequest < 0 || w.CPULimit < 0 {
		return errors.Wrap(ErrInvalidCPU, "limit or request less than 0")
	}
	if w.CPURequest == 0 && w.CPUBind {
		return errors.Wrap(ErrInvalidCPU, "unlimited request with bind")
	}
	if w.MemRequest == 0 && w.MemLimit > 0 {
		w.MemRequest = w.MemLimit
	}
	if w.MemLimit > 0 && w.MemRequest > 0 && w.MemLimit < w.MemRequest {
		w.MemLimit = w.MemRequest
	}
	if w.CPURequest > 0 && w.CPULimit > 0 && w.CPULimit < w.CPURequest {
		w.CPULimit = w.CPURequest
	}
	// if CPUBind=true, set cpu request=limit to solve the dilemma
	// only deal with cpu limit>request but not vice versa
	if w.CPUBind && w.CPURequest > 0 && w.CPULimit > 0 && w.CPULimit > w.CPURequest {
		w.CPURequest = w.CPULimit
	}
	return nil
}

// ParseFromString .
func (w *WorkloadResourceOpts) ParseFromString(str string) (err error) {
	rawParams := RawParams{}
	if err = json.Unmarshal([]byte(str), &rawParams); err != nil {
		return err
	}

	w.KeepCPUBind = rawParams.Bool("keep-cpu-bind")
	w.CPUBind = rawParams.Bool("cpu-bind")
	w.CPURequest = rawParams.Float64("cpu-request")
	w.CPULimit = rawParams.Float64("cpu-limit")
	// check if cpu shortcut is set
	if cpu := rawParams.Float64("cpu"); cpu > 0 {
		w.CPURequest = cpu
		w.CPULimit = cpu
	}
	if w.MemRequest, err = utils.ParseRAMInHuman(rawParams.String("memory-request")); err != nil {
		return err
	}
	if w.MemLimit, err = utils.ParseRAMInHuman(rawParams.String("memory-limit")); err != nil {
		return err
	}
	// check if mem shortcut is set
	if rawParams.IsSet("memory") {
		var mem int64
		if mem, err = utils.ParseRAMInHuman(rawParams.String("memory")); err != nil {
			return err
		}
		w.MemLimit = mem
		w.MemRequest = mem
	}

	return nil
}

// NodeResourceOpts includes all possible fields passed by eru-core for editing node
type NodeResourceOpts struct {
	CPUMap     CPUMap     `json:"cpu_map"`
	Memory     int64      `json:"memory"`
	NUMA       NUMA       `json:"numa"`
	NUMAMemory NUMAMemory `json:"numa_memory"`
	Delta      bool       `json:"delta"`
}

func (n *NodeResourceOpts) ParseFromString(str string) (err error) {
	rawParams := RawParams{}
	if err = json.Unmarshal([]byte(str), &rawParams); err != nil {
		return err
	}

	if n.CPUMap == nil {
		n.CPUMap = CPUMap{}
	}

	if cpu := rawParams.Int64("cpu"); cpu > 0 {
		share := rawParams.Int64("share")
		if share == 0 {
			share = 100
		}

		for i := int64(0); i < cpu; i++ {
			n.CPUMap[fmt.Sprintf("%v", i)] = int(share)
		}
	} else {
		cpuList := rawParams.String("cpu")
		if cpuList != "" {
			cpuMapList := strings.Split(cpuList, ",")
			for _, cpus := range cpuMapList {
				cpuConfigs := strings.Split(cpus, ":")
				pieces, err := strconv.ParseInt(cpuConfigs[1], 10, 32)
				if err != nil {
					return err
				}
				cpuID := cpuConfigs[0]
				n.CPUMap[cpuID] = int(pieces)
			}
		}
	}

	if n.Memory, err = utils.ParseRAMInHuman(rawParams.String("memory")); err != nil {
		return err
	}
	n.NUMA = NUMA{}
	n.NUMAMemory = NUMAMemory{}

	for index, cpuList := range rawParams.StringSlice("numa-cpu") {
		nodeID := fmt.Sprintf("%d", index)
		for _, cpuID := range strings.Split(cpuList, ",") {
			n.NUMA[cpuID] = nodeID
		}
	}

	for index, memoryStr := range rawParams.StringSlice("numa-memory") {
		nodeID := fmt.Sprintf("%d", index)
		mem, err := utils.ParseRAMInHuman(memoryStr)
		if err != nil {
			return err
		}
		n.NUMAMemory[nodeID] = mem
	}

	n.Delta = rawParams.Bool("delta")

	return nil
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
	Remap    bool    `json:"remap"`
}
