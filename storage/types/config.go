package types

import coretypes "github.com/projecteru2/core/types"

// SchedConfig holds scheduler config
type SchedConfig struct {
	MaxDeployCount int `yaml:"max_deploy_count" required:"false" default:"10000"` // max deploy count of each node
}

// Config .
type Config struct {
	ETCD      coretypes.EtcdConfig `yaml:"etcd"`
	Scheduler SchedConfig          `yaml:"scheduler"`
}
