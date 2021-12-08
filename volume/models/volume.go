package models

import (
	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/volume/types"
	"github.com/projecteru2/core/store/etcdv3/meta"
)

// Volume .
type Volume struct {
	config *types.Config
	store meta.KV
}

// NewVolume .
func NewVolume(config *types.Config) (*Volume, error) {
	v := &Volume{config: config}
	var err error
	if len(config.ETCD.Machines) > 0 {
		v.store, err = meta.NewETCD(config.ETCD, nil)
		if err != nil {
			logrus.Errorf("[NewVolume] failed to create etcd client, err: %v", err)
			return nil, err
		}
	}
	return v, nil
}