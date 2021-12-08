package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/volume/types"
)

var allocCommand = &cli.Command{
	Name:   "alloc",
	Usage:  "allocate resources (pure calculation)",
	Action: cmdAlloc,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.IntFlag{
			Name:     "deploy",
			Usage:    "deploy count",
			Required: false,
			Value:    1,
		},
		&cli.StringFlag{
			Name:     "resource-opts",
			Usage:    "resource options",
			Required: true,
		},
	},
}

func cmdAlloc(c *cli.Context) error {
	// volume alloc --node xxx-node --deploy 3 --resource-opts {"cpu": 3, "cpu-bind": true}
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	deployCount := c.Int("deploy")
	resourceOpts := &types.WorkloadResourceOpts{}

	if err := resourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
		logrus.Errorf("[cmdAlloc] invalid resource opts, err: %v", err)
		return err
	}

	engineArgs, resourceArgs, err := volume.Alloc(c.Context, node, deployCount, resourceOpts)
	if err != nil {
		logrus.Errorf("[cmdAlloc] failed to alloc resource, err: %v", err)
		return err
	}

	printResult(map[string]interface{}{
		"engine_args":   engineArgs,
		"resource_args": resourceArgs,
	})

	return nil
}