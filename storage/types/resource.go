package types

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"

	pluginutils "github.com/projecteru2/core-plugins/utils"
)

// WorkloadResourceOpts .
type WorkloadResourceOpts struct {
	StorageRequest int64 `json:"storage-request"`
	StorageLimit   int64 `json:"storage-limit"`
}

// Validate .
func (w *WorkloadResourceOpts) Validate() error {
	if w.StorageLimit < 0 || w.StorageRequest < 0 {
		return errors.Wrap(ErrInvalidStorage, "storage limit or request less than 0")
	}
	if w.StorageLimit > 0 && w.StorageRequest == 0 {
		w.StorageRequest = w.StorageLimit
	}
	if w.StorageLimit > 0 && w.StorageRequest > 0 && w.StorageRequest > w.StorageLimit {
		w.StorageLimit = w.StorageRequest // soft limit storage size
	}
	return nil
}

// ParseFromString .
func (w *WorkloadResourceOpts) ParseFromString(str string) (err error) {
	rawParams := pluginutils.RawParams{}
	if err = json.Unmarshal([]byte(str), &rawParams); err != nil {
		return err
	}

	if w.StorageRequest, err = pluginutils.ParseRAMInHuman(rawParams.String("storage-request")); err != nil {
		return err
	}
	if w.StorageLimit, err = pluginutils.ParseRAMInHuman(rawParams.String("storage-limit")); err != nil {
		return err
	}
	if rawParams.IsSet("storage") {
		if storage, err := pluginutils.ParseRAMInHuman(rawParams.String("storage")); err != nil {
			return err
		} else {
			w.StorageLimit = storage
			w.StorageRequest = storage
		}
	}
	return nil
}

// WorkloadResourceArgs .
type WorkloadResourceArgs struct {
	StorageRequest int64 `json:"storage-request"`
	StorageLimit   int64 `json:"storage-limit"`
}

// VolumeMap .
type VolumeMap map[string]int64

// Total .
func (v VolumeMap) Total() int64 {
	res := int64(0)
	for _, size := range v {
		res += size
	}
	return res
}

// NodeResourceOpts .
type NodeResourceOpts struct {
	Storage int64     `json:"storage"`
	Volumes VolumeMap `json:"volumes"`

	rawParams pluginutils.RawParams
}

// ParseFromString .
func (n *NodeResourceOpts) ParseFromString(str string) (err error) {
	n.rawParams = pluginutils.RawParams{}
	if err = json.Unmarshal([]byte(str), &n.rawParams); err != nil {
		return err
	}

	if n.Storage, err = pluginutils.ParseRAMInHuman(n.rawParams.String("storage")); err != nil {
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
	return nil
}

// SkipEmpty used for setting node resource capacity in absolute mode
func (n *NodeResourceOpts) SkipEmpty(resourceCapacity *NodeResourceArgs) {
	if n == nil {
		return
	}
	if !n.rawParams.IsSet("storage") {
		n.Storage = resourceCapacity.Storage
	}
}

// NodeResourceArgs .
type NodeResourceArgs struct {
	Storage int64 `json:"storage"`
}

// DeepCopy .
func (n *NodeResourceArgs) DeepCopy() *NodeResourceArgs {
	return &NodeResourceArgs{Storage: n.Storage}
}

// Add .
func (n *NodeResourceArgs) Add(n1 *NodeResourceArgs) {
	n.Storage += n1.Storage
}

// Sub .
func (n *NodeResourceArgs) Sub(n1 *NodeResourceArgs) {
	n.Storage -= n1.Storage
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
		n.Usage = &NodeResourceArgs{Storage: 0}
	}

	if n.Capacity.Storage < 0 {
		return errors.Wrap(ErrInvalidStorage, "storage capacity can't be negative")
	}
	if n.Usage.Storage < 0 {
		return errors.Wrap(ErrInvalidStorage, "storage usage can't be negative")
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
	Storage int64 `json:"storage"`
}

// WorkloadResourceArgsMap .
type WorkloadResourceArgsMap map[string]*WorkloadResourceArgs

// ParseFromString .
func (w *WorkloadResourceArgsMap) ParseFromString(str string) error {
	return json.Unmarshal([]byte(str), w)
}
