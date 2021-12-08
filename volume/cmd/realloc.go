package cmd

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/volume/types"
)

var reallocCommand = &cli.Command{
	Name:   "realloc",
	Usage:  "reallocate resources (pure calculation)",
	Action: cmdRealloc,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringFlag{
			Name: "old",
			Usage: "origin resource usage",
			Required: true,
		},
		&cli.StringFlag{
			Name: "resource-opts",
			Usage: "resource options in delta form",
			Required: true,
		},
	},
}

func cmdRealloc(c *cli.Context) error {
	// volume realloc --node xxx-node --old {"cpu": 1} --resource-opts {"cpu": 2}
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	oldResourceArgs := &types.WorkloadResourceArgs{}
	newResourceOpts := &types.WorkloadResourceOpts{}

	if err := json.Unmarshal([]byte(c.String("old")), oldResourceArgs); err != nil {
		logrus.Errorf("[cmdAlloc] invalid old resource args, err: %v", err)
		return err
	}
	if err := newResourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
		logrus.Errorf("[cmdAlloc] invalid new resource opts, err: %v", err)
		return err
	}

	engineArgs, deltaWorkloadResourceArgs, finalWorkloadResourceArgs, err := volume.Realloc(c.Context, node, oldResourceArgs, newResourceOpts)
	if err != nil {
		logrus.Errorf("[cmdRealloc] failed to realloc, err: %v", err)
		return err
	}
	printResult(map[string]interface{}{
		"engine_args":   engineArgs,
		"delta":         deltaWorkloadResourceArgs,
		"resource_args": finalWorkloadResourceArgs,
	})

	return nil
}
