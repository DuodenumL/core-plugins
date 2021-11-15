package cmd

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/cpumem/types"
)

var remapCommand = &cli.Command{
	Name:   "remap",
	Usage:  "remap (pure calculation)",
	Action: cmdRemap,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringFlag{
			Name: "workload-map",
			Usage: "the resource usage of all workloads belonging to this node",
			Required: true,
		},
	},
}

func cmdRemap(c *cli.Context) error {
	// cpumem remap --node xxx-node --workload-map {"workload-1": {"cpu": 1}}
	cpumem, err := newCPUMem(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	workloadMap := map[string]*types.WorkloadResourceArgs{}

	if err := json.Unmarshal([]byte(c.String("workload-map")), &workloadMap); err != nil {
		logrus.Errorf("[cmdRemap] invalid workload map, err: %v", err)
		return err
	}

	engineArgsMap, err := cpumem.Remap(c.Context, node, workloadMap)
	if err != nil {
		logrus.Errorf("[cmdRemap] failed to remap, err: %v", err)
		return err
	}

	printResult(map[string]interface{}{
		"engine_args_map": engineArgsMap,
	})

	return nil
}
