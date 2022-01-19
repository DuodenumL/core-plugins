package models

import (
	"github.com/sirupsen/logrus"

	"github.com/projecteru2/core-plugins/storage/types"
	"github.com/projecteru2/core/store/etcdv3/meta"
)

// Storage .
type Storage struct {
	config *types.Config
	store  meta.KV
}

// NewStorage .
func NewStorage(config *types.Config) (*Storage, error) {
	s := &Storage{config: config}
	var err error
	if len(config.ETCD.Machines) > 0 {
		s.store, err = meta.NewETCD(config.ETCD, nil)
		if err != nil {
			logrus.Errorf("[NewStorage] failed to create etcd client, err: %v", err)
			return nil, err
		}
	}
	return s, nil
}
