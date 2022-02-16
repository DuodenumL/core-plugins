package models

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/volume/schedule"
	"github.com/projecteru2/core-plugins/volume/types"
)

// GetDeployArgs .
func (v *Volume) GetDeployArgs(ctx context.Context, node string, deployCount int, opts *types.WorkloadResourceOpts) ([]*types.EngineArgs, []*types.WorkloadResourceArgs, error) {
	if err := opts.Validate(); err != nil {
		logrus.Errorf("[GetDeployArgs] invalid resource opts %+v, err: %v", opts, err)
		return nil, nil, err
	}

	resourceInfo, err := v.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[GetDeployArgs] failed to get resource info of node %v, err: %v", node, err)
		return nil, nil, err
	}

	return v.doAlloc(resourceInfo, deployCount, opts)
}

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func getVolumePlanLimit(bindings types.VolumeBindings, volumePlan types.VolumePlan) types.VolumePlan {
	volumeBindingToVolumeMap := map[[3]string]types.VolumeMap{}
	for binding, volumeMap := range volumePlan {
		volumeBindingToVolumeMap[binding.GetMapKey()] = volumeMap
	}

	volumePlanLimit := types.VolumePlan{}

	for _, binding := range bindings {
		if volumeMap, ok := volumeBindingToVolumeMap[binding.GetMapKey()]; ok {
			volumePlanLimit[binding] = types.VolumeMap{volumeMap.GetDevice(): maxInt64(binding.SizeInBytes, volumeMap.GetSize())}
		}
	}
	return volumePlanLimit
}

func (v *Volume) doAlloc(resourceInfo *types.NodeResourceInfo, deployCount int, opts *types.WorkloadResourceOpts) ([]*types.EngineArgs, []*types.WorkloadResourceArgs, error) {
	volumePlans := schedule.GetVolumePlans(resourceInfo, opts.VolumesRequest, v.config.Scheduler.MaxDeployCount)
	if opts.StorageRequest > 0 {
		storageCapacity := int((resourceInfo.Capacity.Storage - resourceInfo.Usage.Storage) / opts.StorageRequest)
		if storageCapacity < len(volumePlans) {
			volumePlans = volumePlans[:storageCapacity]
		}
	}

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
		engineArgs := &types.EngineArgs{Storage: opts.StorageLimit}
		for _, binding := range opts.VolumesLimit.ApplyPlan(volumePlan) {
			engineArgs.Volumes = append(engineArgs.Volumes, binding.ToString(true))
		}

		resourceArgs := &types.WorkloadResourceArgs{
			VolumesRequest:    opts.VolumesRequest,
			VolumesLimit:      opts.VolumesLimit,
			VolumePlanRequest: volumePlan,
			VolumePlanLimit:   getVolumePlanLimit(opts.VolumesLimit, volumePlan),
			StorageRequest:    opts.StorageRequest,
			StorageLimit:      opts.StorageLimit,
		}

		resEngineArgs = append(resEngineArgs, engineArgs)
		resResourceArgs = append(resResourceArgs, resourceArgs)
	}

	return resEngineArgs, resResourceArgs, nil
}
