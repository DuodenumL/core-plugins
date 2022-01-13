package command

import (
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/volume/models"
	"github.com/projecteru2/core-plugins/volume/types"
)

var VolumeCommands = []*cli.Command{
	generateGetDeployArgsCommand(newVolume, &types.WorkloadResourceOpts{}),
	generateGetMostIdleNodeCommand(newVolume),
	generateSetNodeResourceInfoCommand(newVolume, &types.WorkloadResourceArgs{}, &types.WorkloadResourceArgs{}),
	generateSetNodeResourceUsageCommand(newVolume, &types.NodeResourceOpts{}, &types.NodeResourceArgs{}, &types.WorkloadResourceArgs{}),
	generateSetNodeResourceCapacityCommand(newVolume, &types.NodeResourceOpts{}, &types.NodeResourceArgs{}),
	generateGetNodeResourceInfoCommand(newVolume, &types.WorkloadResourceArgsMap{}),
	generateGetReallocArgsCommand(newVolume, &types.WorkloadResourceArgs{}, &types.WorkloadResourceOpts{}),
	generateGetRemapArgsCommand(newVolume, &types.WorkloadResourceArgsMap{}),
	generateAddNodeCommand(newVolume, &types.NodeResourceOpts{}),
	generateRemoveNodeCommand(newVolume),
	generateGetNodesDeployCapacityCommand(newVolume, &types.WorkloadResourceOpts{}),
}

func newVolume(c *cli.Context) (interface{}, error) {
	configPath := c.String("config")
	config := &types.Config{}
	if err := configor.Load(config, configPath); err != nil {
		logrus.Errorf("[newVolume] failed to load config, err: %v", err)
		return nil, err
	}

	volume, err := models.NewVolume(config)
	if err != nil {
		logrus.Errorf("[newVolume] failed to init volume, err: %v", err)
		return nil, err
	}
	return volume, nil
}
