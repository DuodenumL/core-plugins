package command

import (
	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/cpumem/models"
	"github.com/projecteru2/core-plugins/cpumem/types"
)

var CPUMemCommands = []*cli.Command{
	generateGetDeployArgsCommand(newCPUMem, &types.WorkloadResourceOpts{}),
	generateGetMostIdleNodeCommand(newCPUMem),
	generateSetNodeResourceInfoCommand(newCPUMem, &types.WorkloadResourceArgs{}, &types.WorkloadResourceArgs{}),
	generateSetNodeResourceUsageCommand(newCPUMem, &types.NodeResourceOpts{}, &types.NodeResourceArgs{}, &types.WorkloadResourceArgs{}),
	generateSetNodeResourceCapacityCommand(newCPUMem, &types.NodeResourceOpts{}, &types.NodeResourceArgs{}),
	generateGetNodeResourceInfoCommand(newCPUMem, &types.WorkloadResourceArgsMap{}),
	generateGetReallocArgsCommand(newCPUMem, &types.WorkloadResourceArgs{}, &types.WorkloadResourceOpts{}),
	generateGetRemapArgsCommand(newCPUMem, &types.WorkloadResourceArgsMap{}),
	generateAddNodeCommand(newCPUMem, &types.NodeResourceOpts{}),
	generateRemoveNodeCommand(newCPUMem),
	generateGetNodesDeployCapacityCommand(newCPUMem, &types.WorkloadResourceOpts{}),
}

func newCPUMem(c *cli.Context) (interface{}, error) {
	configPath := c.String("config")
	config := &types.Config{}
	if err := configor.Load(config, configPath); err != nil {
		logrus.Errorf("[newCPUMem] failed to load config, err: %v", err)
		return nil, err
	}

	cpuMem, err := models.NewCPUMem(config)
	if err != nil {
		logrus.Errorf("[newCPUMem] failed to init cpumem, err: %v", err)
		return nil, err
	}
	return cpuMem, nil
}
