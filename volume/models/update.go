package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/volume/types"
)

// UpdateNodeResourceUsage .
func (v *Volume) UpdateNodeResourceUsage(ctx context.Context, node string, resourceArgsList []*types.WorkloadResourceArgs, incr bool) error {
	resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[UpdateNodeResourceUsage] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	for _, resourceArgs := range resourceArgsList {
		for _, volumeMap := range resourceArgs.VolumePlanRequest {
			if incr {
				resourceInfo.Usage.Volumes.Add(volumeMap)
			} else {
				resourceInfo.Usage.Volumes.Sub(volumeMap)
			}
		}
	}

	return v.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

// UpdateNodeResourceCapacity .
func (v *Volume) UpdateNodeResourceCapacity(ctx context.Context, node string, resourceOpts *types.NodeResourceOpts, incr bool) error {
	resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[UpdateNodeResourceCapacity] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	if incr {
		resourceInfo.Capacity.Volumes.Add(resourceOpts.Volumes)
	} else {
		resourceInfo.Capacity.Volumes.Sub(resourceOpts.Volumes)
	}

	return v.doSetNodeResourceInfo(ctx, node, resourceInfo)
}
