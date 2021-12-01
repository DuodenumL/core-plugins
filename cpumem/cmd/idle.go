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
	// cpumem get-idle --node xxx --node xxx
	cpumem, err := newCPUMem(c)
	if err != nil {
		return err
	}

	nodes := c.StringSlice("node")

	node, priority, err := cpumem.GetMostIdleNode(c.Context, nodes)
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
