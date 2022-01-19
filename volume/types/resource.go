package types

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/pkg/errors"

	pluginutils "github.com/projecteru2/core-plugins/utils"
)

// WorkloadResourceOpts .
type WorkloadResourceOpts struct {
	VolumesRequest VolumeBindings `json:"volumes_request"`
	VolumesLimit   VolumeBindings `json:"volumes_limit"`
}

// Validate .
func (w *WorkloadResourceOpts) Validate() error {
	if len(w.VolumesRequest) == 0 && len(w.VolumesLimit) == 0 {
		return nil
	}
	if len(w.VolumesLimit) > 0 && len(w.VolumesRequest) == 0 {
		w.VolumesRequest = w.VolumesLimit
	}
	if len(w.VolumesRequest) != len(w.VolumesLimit) {
		return errors.Wrap(ErrInvalidVolume, "different length of request and limit")
	}

	sortFunc := func(volumeBindings []*VolumeBinding) func(i, j int) bool {
		return func(i, j int) bool {
			return volumeBindings[i].ToString(false) < volumeBindings[j].ToString(false)
		}
	}

	sort.Slice(w.VolumesRequest, sortFunc(w.VolumesRequest))
	sort.Slice(w.VolumesLimit, sortFunc(w.VolumesLimit))

	for i := range w.VolumesRequest {
		request := w.VolumesRequest[i]
		limit := w.VolumesLimit[i]
		if request.Source != limit.Source || request.Destination != limit.Destination || request.Flags != limit.Flags {
			return errors.Wrap(ErrInvalidVolume, "request and limit not match")
		}
		if request.SizeInBytes > 0 && limit.SizeInBytes > 0 && request.SizeInBytes > limit.SizeInBytes {
			limit.SizeInBytes = request.SizeInBytes
		}
	}
	return nil
}

// ParseFromString .
func (w *WorkloadResourceOpts) ParseFromString(str string) (err error) {
	rawParams := pluginutils.RawParams{}
	if err = json.Unmarshal([]byte(str), &rawParams); err != nil {
		return err
	}

	if w.VolumesRequest, err = NewVolumeBindings(rawParams.OneOfStringSlice("volumes-request", "volume-request")); err != nil {
		return err
	}
	if w.VolumesLimit, err = NewVolumeBindings(rawParams.OneOfStringSlice("volumes", "volume", "volume-limit")); err != nil {
		return err
	}
	return nil
}

// WorkloadResourceArgs .
type WorkloadResourceArgs struct {
	VolumesRequest VolumeBindings `json:"volumes_request"`
	VolumesLimit   VolumeBindings `json:"volumes_limit"`

	VolumePlanRequest VolumePlan `json:"volume_plan_request"`
	VolumePlanLimit   VolumePlan `json:"volume_plan_limit"`
}

// NodeResourceOpts .
type NodeResourceOpts struct {
	Volumes VolumeMap `json:"volumes"`

	rawParams pluginutils.RawParams
}

// ParseFromString .
func (n *NodeResourceOpts) ParseFromString(str string) (err error) {
	n.rawParams = pluginutils.RawParams{}
	if err = json.Unmarshal([]byte(str), &n.rawParams); err != nil {
		return err
	}

	volumes := VolumeMap{}
	for _, volume := range n.rawParams.StringSlice("volumes") {
		parts := strings.Split(volume, ":")
		if len(parts) != 2 {
			return errors.Wrap(ErrInvalidVolume, "volume should have 2 parts")
		}

		capacity, err := pluginutils.ParseRAMInHuman(parts[1])
		if err != nil {
			return err
		}
		volumes[parts[0]] = capacity
	}
	n.Volumes = volumes
	return
}

// SkipEmpty used for setting node resource capacity in absolute mode
func (n *NodeResourceOpts) SkipEmpty(resourceCapacity *NodeResourceArgs) {
	if n == nil {
		return
	}
	if !n.rawParams.IsSet("volumes") {
		n.Volumes = resourceCapacity.Volumes
	}
}

// NodeResourceArgs .
type NodeResourceArgs struct {
	Volumes VolumeMap `json:"volumes"`
}

// DeepCopy .
func (n *NodeResourceArgs) DeepCopy() *NodeResourceArgs {
	return &NodeResourceArgs{Volumes: n.Volumes.DeepCopy()}
}

// RemoveEmpty .
func (n *NodeResourceArgs) RemoveEmpty(n1 *NodeResourceArgs) {
	for device, size := range n1.Volumes {
		if size == 0 {
			delete(n.Volumes, device)
		}
	}
}

// Add .
func (n *NodeResourceArgs) Add(n1 *NodeResourceArgs) {
	for k, v := range n1.Volumes {
		n.Volumes[k] += v
	}
}

// Sub .
func (n *NodeResourceArgs) Sub(n1 *NodeResourceArgs) {
	for k, v := range n1.Volumes {
		n.Volumes[k] -= v
	}
}

// NodeResourceInfo .
type NodeResourceInfo struct {
	Capacity *NodeResourceArgs
	Usage    *NodeResourceArgs
}

// Validate .
func (n *NodeResourceInfo) Validate() error {
	if n.Capacity == nil {
		return ErrInvalidCapacity
	}
	if n.Usage == nil {
		n.Usage = &NodeResourceArgs{Volumes: VolumeMap{}}
		for device := range n.Capacity.Volumes {
			n.Usage.Volumes[device] = 0
		}
	}

	for key, value := range n.Capacity.Volumes {
		if value < 0 {
			return errors.Wrap(ErrInvalidVolume, "volume size should not be less than 0")
		}
		if usage, ok := n.Usage.Volumes[key]; ok && (usage > value || usage < 0) {
			return errors.Wrap(ErrInvalidVolume, "invalid size in usage")
		}
	}
	for key := range n.Usage.Volumes {
		if _, ok := n.Usage.Volumes[key]; !ok {
			return errors.Wrap(ErrInvalidVolume, "invalid key in usage")
		}
	}
	return nil
}

// GetAvailableResource .
func (n *NodeResourceInfo) GetAvailableResource() *NodeResourceArgs {
	res := n.Capacity.DeepCopy()
	res.Sub(n.Usage)
	return res
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
	Volumes       []string `json:"volumes"`
	VolumeChanged bool     `json:"volume_changed"` // indicates whether the realloc request includes new volumes
}

// WorkloadResourceArgsMap .
type WorkloadResourceArgsMap map[string]*WorkloadResourceArgs

// ParseFromString .
func (w *WorkloadResourceArgsMap) ParseFromString(str string) error {
	return json.Unmarshal([]byte(str), w)
}
