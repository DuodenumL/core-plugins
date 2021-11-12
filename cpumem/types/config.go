package types

import coretypes "github.com/projecteru2/core/types"

// LogConfig .
type LogConfig struct {
	Stderr bool
}

// SchedConfig holds scheduler config
type SchedConfig struct {
	MaxShare  int `yaml:"maxshare" required:"true" default:"-1"`   // comlpex scheduler use maxshare
	ShareBase int `yaml:"sharebase" required:"true" default:"100"` // how many pieces for one core
}

// Config .
type Config struct {
	ETCD      coretypes.EtcdConfig `yaml:"etcd"`
	Log       LogConfig            `yaml:"log"`
	Scheduler SchedConfig          `yaml:"scheduler"`
}
