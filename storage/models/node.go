package models

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/storage/types"
	coretypes "github.com/projecteru2/core/types"
)

// AddNode .
func (s *Storage) AddNode(ctx context.Context, node string, resourceOpts *types.NodeResourceOpts) (*types.NodeResourceInfo, error) {
	if _, err := s.doGetNodeResourceInfo(ctx, node); err != nil {
		if !errors.Is(err, coretypes.ErrBadCount) {
			logrus.Errorf("[AddNode] failed to get resource info of node %v, err: %v", node, err)
			return nil, err
		}
	} else {
		return nil, types.ErrNodeExists
	}

	resourceInfo := &types.NodeResourceInfo{
		Capacity: &types.NodeResourceArgs{
			Storage: resourceOpts.Storage + resourceOpts.Volumes.Total(),
		},
		Usage: nil,
	}

	return resourceInfo, s.doSetNodeResourceInfo(ctx, node, resourceInfo)
}

// RemoveNode .
func (s *Storage) RemoveNode(ctx context.Context, node string) error {
	if _, err := s.store.Delete(ctx, fmt.Sprintf(NodeResourceInfoKey, node)); err != nil {
		logrus.Errorf("[doSetNodeResourceInfo] faield to delete node %v, err: %v", node, err)
		return err
	}
	return nil
}
