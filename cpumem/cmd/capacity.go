package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/cpumem/types"
)

var capacityCommand = &cli.Command{
	Name:   "get-capacity",
	Usage:  "calculate how many workloads can be allocated according to the specified resource options (pure calculation)",
	Action: cmdCapacity,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "resource-opts",
			Usage:    "resource options",
			Required: true,
		},
	},
}

func cmdCapacity(c *cli.Context) error {
	// cpumem get-capacity --node xxx-node --node xxx-node --resource-opts {"cpu": 1}
	cpumem, err := newCPUMem(c)
	if err != nil {
		return err
	}

	nodes := c.StringSlice("node")
	resourceOpts := &types.WorkloadResourceOpts{}

	if err := resourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
		logrus.Errorf("[cmdCapacity] invalid resource opts, err: %v", err)
		return err
	}

	nodeCapacityInfoMap, total, err := cpumem.GetNodesCapacity(c.Context, nodes, resourceOpts)
	if err != nil {
		logrus.Errorf("[cmdCapacity] failed to get capacity, err: %v", err)
		return err
	}
	result := map[string]interface{}{
		"nodes": nodeCapacityInfoMap,
		"total": total,
	}
	printResult(result)

	return nil
}
