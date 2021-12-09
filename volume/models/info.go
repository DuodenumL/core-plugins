package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sanity-io/litter"
	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/volume/types"
)

const NodeResourceInfoKey = "/resource/volume/%s"

// GetNodeResourceInfo .
func (v *Volume) GetNodeResourceInfo(ctx context.Context, node string, workloadResourceMap map[string]*types.WorkloadResourceArgs, fix bool) (*types.NodeResourceInfo, []string, error) {
	resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		return nil, nil, err
	}

	diffs := []string{}

	totalVolumeMap := types.VolumeMap{}

	for _, args := range workloadResourceMap {
		for _, volumeMap := range args.VolumePlanRequest {
			totalVolumeMap.Add(volumeMap)
		}
	}

	for volume, size := range totalVolumeMap {
		if resourceInfo.Usage.Volumes[volume] != size {
			diffs = append(diffs, fmt.Sprintf("node.Volumes[%v] != sum(workload.Volumes[%v]: %v != %v)", volume, volume, resourceInfo.Usage.Volumes[volume], size))
		}
	}

	if fix {
		resourceInfo.Usage = &types.NodeResourceArgs{
			Volumes: totalVolumeMap,
		}
		if err = v.doSetNodeResourceInfo(ctx, node, resourceInfo); err != nil {
			logrus.Warnf("[GetNodeResourceInfo] failed to fix node resource, err: %v", err)
			diffs = append(diffs, "fix failed")
		}
	}

	return resourceInfo, diffs, nil
}

// SetNodeResourceInfo .
func (v *Volume) SetNodeResourceInfo(ctx context.Context, node string, resourceCapacity *types.NodeResourceArgs, resourceUsage *types.NodeResourceArgs) error {
	resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[SetNodeResourceInfo] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	resourceInfo.Capacity = resourceCapacity
	resourceInfo.Usage = resourceUsage

	return v.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

func (v *Volume) doGetNodeResourceInfo(ctx context.Context, node string) (*types.NodeResourceInfo, error) {
	resourceInfo := &types.NodeResourceInfo{}
	resp, err := v.store.GetOne(ctx, fmt.Sprintf(NodeResourceInfoKey, node))
	if err != nil {
		logrus.Errorf("[doGetNodeResourceInfo] failed to get node resource info of node %v, err: %v", node, err)
		return nil, err
	}
	if err = json.Unmarshal(resp.Value, resourceInfo); err != nil {
		logrus.Errorf("[doGetNodeResourceInfo] failed to unmarshal node resource info of node %v, err: %v", node, err)
		return nil, err
	}
	return resourceInfo, nil
}

func (v *Volume) doSetNodeResourceInfo(ctx context.Context, node string, resourceInfo *types.NodeResourceInfo) error {
	if err := resourceInfo.Validate(); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] invalid resource info %v, err: %v", litter.Sdump(resourceInfo), err)
		return err
	}

	data, err := json.Marshal(resourceInfo)
	if err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to marshal resource info %+v, err: %v", resourceInfo, err)
		return err
	}

	if _, err = v.store.Put(ctx, fmt.Sprintf(NodeResourceInfoKey, node), string(data)); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to put resource info %+v, err: %v", resourceInfo, err)
		return err
	}
	return nil
}
