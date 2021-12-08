package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/volume/schedule"
	"github.com/projecteru2/core-plugins/volume/types"
)

// Alloc .
func (v *Volume) Alloc(ctx context.Context, node string, deployCount int, opts *types.WorkloadResourceOpts) ([]*types.EngineArgs, []*types.WorkloadResourceArgs, error) {
	if err := opts.Validate(); err != nil {
		logrus.Errorf("[Alloc] invalid resource opts %+v, err: %v", opts, err)
		return nil, nil, err
	}

	resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[Alloc] failed to get resource info of node %v, err: %v", node, err)
		return nil, nil, err
	}

	return v.doAlloc(resourceInfo, deployCount, opts)
}

func (v *Volume) doAlloc(resourceInfo *types.NodeResourceInfo, deployCount int, opts *types.WorkloadResourceOpts) ([]*types.EngineArgs, []*types.WorkloadResourceArgs, error) {
	volumePlans := schedule.GetVolumePlans(resourceInfo, opts.VolumesRequest, nil, v.config.Scheduler.MaxDeployCount)
	if len(volumePlans) < deployCount {
		return nil, nil, types.ErrInsufficientResource
	}

	volumePlans = volumePlans[:deployCount]
	resEngineArgs := []*types.EngineArgs{}
	resResourceArgs := []*types.WorkloadResourceArgs{}

	volumeSizeLimitMap := map[*types.VolumeBinding]int64{}
	for _, binding := range opts.VolumesLimit {
		volumeSizeLimitMap[binding] = binding.SizeInBytes
	}

	for _, volumePlan := range volumePlans {
		engineArgs := &types.EngineArgs{}
		for _, binding := range opts.VolumesLimit {
			engineArgs.Volumes = append(engineArgs.Volumes, binding.ToString(true))
		}

		volumePlanLimit := types.VolumePlan{}
		for binding, volumeMap := range volumePlan {
			for device := range volumeMap {
				volumePlanLimit[binding] = types.VolumeMap{device: volumeSizeLimitMap[binding]}
			}
		}

		resourceArgs := &types.WorkloadResourceArgs{
			VolumesRequest:    opts.VolumesRequest,
			VolumesLimit:      opts.VolumesLimit,
			VolumePlanRequest: volumePlan,
			VolumePlanLimit:   volumePlanLimit,
		}

		resEngineArgs = append(resEngineArgs, engineArgs)
		resResourceArgs = append(resResourceArgs, resourceArgs)
	}

	return resEngineArgs, resResourceArgs, nil
}
