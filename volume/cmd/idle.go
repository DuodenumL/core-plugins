package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var getIdleCommand = &cli.Command{
	Name:   "get-idle",
	Usage:  "get most idle nodes",
	Action: cmdGetIdle,
	Flags: []cli.Flag{
		&cli.StringSliceFlag{
			Name:     "node",
			Usage:    "node names",
			Required: true,
		},
	},
}

func cmdGetIdle(c *cli.Context) error {
	// volume get-idle --node xxx --node xxx
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	nodes := c.StringSlice("node")

	node, priority, err := volume.GetMostIdleNode(c.Context, nodes)
	if err != nil {
		logrus.Errorf("[cmdGetIdle] failed to get the most idle node, err: %v", err)
		return err
	}

	printResult(map[string]interface{}{
		"node": node,
		"priority": priority,
	})
	return nil
}
