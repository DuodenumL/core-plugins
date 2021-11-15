package cmd

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/cpumem/types"
)

var updateCapacityCommand = &cli.Command{
	Name:   "update-capacity",
	Usage:  "update node resource capacity",
	Action: cmdUpdateNodeResourceCapacity,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "resource-opts",
			Usage:    "resource options",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "decr",
			Usage:    "decrease",
			Value:    false,
			Required: false,
		},
	},
}

var updateUsageCommand = &cli.Command{
	Name:   "update-usage",
	Usage:  "update node resource usage",
	Action: cmdUpdateNodeResourceUsage,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringSliceFlag{
			Name:     "resource-args",
			Usage:    "resource args",
			Required: true,
		},
		&cli.BoolFlag{
			Name:     "decr",
			Usage:    "decrease",
			Value:    false,
			Required: false,
		},
	},
}

func cmdUpdateNodeResourceCapacity(c *cli.Context) error {
	// cpumem update-capacity --node xxx-node --resource-opts {"cpu": 1} --incr
	cpumem, err := newCPUMem(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	incr := !c.Bool("decr")
	resourceOpts := &types.NodeResourceOpts{}

	if err := resourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
		logrus.Errorf("[cmdUpdateNodeResourceCapacity] invalid resource opts, err: %v", err)
		return err
	}

	if err = cpumem.UpdateNodeResourceCapacity(c.Context, node, resourceOpts, incr); err != nil {
		logrus.Errorf("[cmdUpdateNodeResourceCapacity] failed to update node resource capacity, err: %v", err)
		return err
	}

	return nil
}

func cmdUpdateNodeResourceUsage(c *cli.Context) error {
	// cpumem update-usage --node xxx-node --resource-args {"cpu": 1} --incr
	cpumem, err := newCPUMem(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	incr := !c.Bool("decr")
	resourceArgsList := []*types.WorkloadResourceArgs{}

	for _, args := range c.StringSlice("resource-args") {
		resourceArgs := &types.WorkloadResourceArgs{}
		if err := json.Unmarshal([]byte(args), &resourceArgs); err != nil {
			logrus.Errorf("[cmdUpdateNodeResourceUsage] invalid resource args, err: %v", err)
			return err
		}
		resourceArgsList = append(resourceArgsList, resourceArgs)
	}

	if err = cpumem.UpdateNodeResourceUsage(c.Context, node, resourceArgsList, incr); err != nil {
		logrus.Errorf("[cmdUpdateNodeResourceUsage] failed to update node resource usage, err: %v", err)
		return err
	}
	return nil
}
