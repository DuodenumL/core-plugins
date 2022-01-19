package command

import (
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/storage/models"
	"github.com/projecteru2/core-plugins/storage/types"
)

var StorageCommands = []*cli.Command{
	generateGetDeployArgsCommand(newStorage, &types.WorkloadResourceOpts{}),
	generateGetMostIdleNodeCommand(newStorage),
	generateSetNodeResourceInfoCommand(newStorage, &types.WorkloadResourceArgs{}, &types.WorkloadResourceArgs{}),
	generateSetNodeResourceUsageCommand(newStorage, &types.NodeResourceOpts{}, &types.NodeResourceArgs{}, &types.WorkloadResourceArgs{}),
	generateSetNodeResourceCapacityCommand(newStorage, &types.NodeResourceOpts{}, &types.NodeResourceArgs{}),
	generateGetNodeResourceInfoCommand(newStorage, &types.WorkloadResourceArgsMap{}),
	generateGetReallocArgsCommand(newStorage, &types.WorkloadResourceArgs{}, &types.WorkloadResourceOpts{}),
	generateGetRemapArgsCommand(newStorage, &types.WorkloadResourceArgsMap{}),
	generateAddNodeCommand(newStorage, &types.NodeResourceOpts{}),
	generateRemoveNodeCommand(newStorage),
	generateGetNodesDeployCapacityCommand(newStorage, &types.WorkloadResourceOpts{}),
}

func newStorage(c *cli.Context) (interface{}, error) {
	configPath := c.String("config")
	config := &types.Config{}
	if err := configor.Load(config, configPath); err != nil {
		logrus.Errorf("[newStorage] failed to load config, err: %v", err)
		return nil, err
	}

	storage, err := models.NewStorage(config)
	if err != nil {
		logrus.Errorf("[newStorage] failed to init storage, err: %v", err)
		return nil, err
	}
	return storage, nil
}
