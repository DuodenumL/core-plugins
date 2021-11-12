package models

import (
	"github.com/projecteru2/core-plugins/cpumem/types"
	"github.com/projecteru2/core/store/etcdv3/meta"
	"github.com/sirupsen/logrus"
)

// CPUMem cpumem plugin
type CPUMem struct {
	config *types.Config
	store meta.KV
}

// NewCPUMem .
func NewCPUMem(config *types.Config) (*CPUMem, error) {
	c := &CPUMem{config: config}
	var err error
	if len(config.ETCD.Machines) > 0 {
		c.store, err = meta.NewETCD(config.ETCD, nil)
		if err != nil {
			logrus.Errorf("[NewCPUMem] failed to create etcd client, err: %v", err)
			return nil, err
		}
	}
	return c, nil
}