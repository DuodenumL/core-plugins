package models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/sanity-io/litter"
	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/storage/types"
)

const NodeResourceInfoKey = "/resource/storage/%s"

// GetNodeResourceInfo .
func (s *Storage) GetNodeResourceInfo(ctx context.Context, node string, workloadResourceMap *types.WorkloadResourceArgsMap, fix bool) (*types.NodeResourceInfo, []string, error) {
	resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		return nil, nil, err
	}

	diffs := []string{}

	totalStorageUsage := int64(0)

	for _, args := range *workloadResourceMap {
		totalStorageUsage += args.StorageRequest
	}

	if resourceInfo.Usage.Storage != totalStorageUsage {
		diffs = append(diffs, fmt.Sprintf("node.Storage != sum(workload.Storage): %v != %v", resourceInfo.Usage.Storage, totalStorageUsage))
	}

	if fix {
		resourceInfo.Usage = &types.NodeResourceArgs{
			Storage: totalStorageUsage,
		}
		if err = s.doSetNodeResourceInfo(ctx, node, resourceInfo); err != nil {
			logrus.Warnf("[GetNodeResourceInfo] failed to fix node resource, err: %v", err)
			diffs = append(diffs, "fix failed")
		}
	}

	return resourceInfo, diffs, nil
}

// priority: node resource opts > node resource args > workload resource args list
func (s *Storage) calculateNodeResourceArgs(origin *types.NodeResourceArgs, nodeResourceOpts *types.NodeResourceOpts, nodeResourceArgs *types.NodeResourceArgs, workloadResourceArgs []*types.WorkloadResourceArgs, delta bool, incr bool) (res *types.NodeResourceArgs) {
	if origin == nil || !delta {
		res = (&types.NodeResourceArgs{}).DeepCopy()
	} else {
		res = origin.DeepCopy()
	}

	if nodeResourceOpts != nil {
		nodeResourceArgs := &types.NodeResourceArgs{
			Storage: nodeResourceOpts.Storage + nodeResourceOpts.Volumes.Total(),
		}

		if incr {
			res.Add(nodeResourceArgs)
		} else {
			res.Sub(nodeResourceArgs)
		}
		return res
	}

	if nodeResourceArgs != nil {
		if incr {
			res.Add(nodeResourceArgs)
		} else {
			res.Sub(nodeResourceArgs)
		}
		return res
	}

	for _, args := range workloadResourceArgs {
		nodeResourceArgs := &types.NodeResourceArgs{
			Storage: args.StorageRequest,
		}
		if incr {
			res.Add(nodeResourceArgs)
		} else {
			res.Sub(nodeResourceArgs)
		}
	}
	return res
}

// SetNodeResourceUsage .
func (s *Storage) SetNodeResourceUsage(ctx context.Context, node string, nodeResourceOpts *types.NodeResourceOpts, nodeResourceArgs *types.NodeResourceArgs, workloadResourceArgs []*types.WorkloadResourceArgs, delta bool, incr bool) (before *types.NodeResourceArgs, after *types.NodeResourceArgs, err error) {
	resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[SetNodeResourceInfo] failed to get resource info of node %v, err: %v", node, err)
		return nil, nil, err
	}

	before = resourceInfo.Usage.DeepCopy()
	resourceInfo.Usage = s.calculateNodeResourceArgs(resourceInfo.Usage, nodeResourceOpts, nodeResourceArgs, workloadResourceArgs, delta, incr)

	if err := s.doSetNodeResourceInfo(ctx, node, resourceInfo); err != nil {
		return nil, nil, err
	}
	return before, resourceInfo.Usage, nil
}

// SetNodeResourceCapacity .
func (s *Storage) SetNodeResourceCapacity(ctx context.Context, node string, nodeResourceOpts *types.NodeResourceOpts, nodeResourceArgs *types.NodeResourceArgs, delta bool, incr bool) (before *types.NodeResourceArgs, after *types.NodeResourceArgs, err error) {
	resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[SetNodeResourceInfo] failed to get resource info of node %v, err: %v", node, err)
		return nil, nil, err
	}

	before = resourceInfo.Capacity.DeepCopy()
	if !delta {
		nodeResourceOpts.SkipEmpty(resourceInfo.Capacity)
	}

	resourceInfo.Capacity = s.calculateNodeResourceArgs(resourceInfo.Usage, nodeResourceOpts, nodeResourceArgs, nil, delta, incr)

	if err := s.doSetNodeResourceInfo(ctx, node, resourceInfo); err != nil {
		return nil, nil, err
	}
	return before, resourceInfo.Capacity, nil
}

// SetNodeResourceInfo .
func (s *Storage) SetNodeResourceInfo(ctx context.Context, node string, resourceCapacity *types.NodeResourceArgs, resourceUsage *types.NodeResourceArgs) error {
	resourceInfo, err := s.doGetNodeResourceInfo(ctx, node)
	if err != nil {
		logrus.Errorf("[SetNodeResourceInfo] failed to get resource info of node %v, err: %v", node, err)
		return err
	}

	resourceInfo.Capacity = resourceCapacity
	resourceInfo.Usage = resourceUsage

	return s.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

func (s *Storage) doGetNodeResourceInfo(ctx context.Context, node string) (*types.NodeResourceInfo, error) {
	resourceInfo := &types.NodeResourceInfo{}
	resp, err := s.store.GetOne(ctx, fmt.Sprintf(NodeResourceInfoKey, node))
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

func (s *Storage) doSetNodeResourceInfo(ctx context.Context, node string, resourceInfo *types.NodeResourceInfo) error {
	if err := resourceInfo.Validate(); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] invalid resource info %v, err: %v", litter.Sdump(resourceInfo), err)
		return err
	}

	data, err := json.Marshal(resourceInfo)
	if err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to marshal resource info %+v, err: %v", resourceInfo, err)
		return err
	}

	if _, err = s.store.Put(ctx, fmt.Sprintf(NodeResourceInfoKey, node), string(data)); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to put resource info %+v, err: %v", resourceInfo, err)
		return err
	}
	return nil
}
