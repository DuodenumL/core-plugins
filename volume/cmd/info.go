package cmd

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/volume/types"
)

var setNodeCommand = &cli.Command{
	Name:   "set-node",
	Usage:  "set the total amount and usage of node resources",
	Action: cmdSetNode,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringFlag{
			Name: "capacity",
			Usage: "resource capacity",
			Required: true,
		},
		&cli.StringFlag{
			Name: "usage",
			Usage: "resource usage",
			Required: true,
		},
	},
}

var getNodeCommand = &cli.Command{
	Name:   "get-node",
	Usage:  "get the total amount and usage of node resources",
	Action: cmdGetNode,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
		&cli.StringFlag{
			Name: "workload-map",
			Usage: "the resource usage of all workloads belonging to this node",
			Value: "{}",
			Required: false,
		},
		&cli.BoolFlag{
			Name: "fix",
			Usage: "fix the resource usage of node according to the resource usage of workload",
			Required: false,
			Value: false,
		},
	},
}

func cmdSetNode(c *cli.Context) error {
	// volume set-node --node xxx-node --capacity {"cpu": 100} --usage {"cpu": 20}
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	resourceCapacity := &types.NodeResourceArgs{}
	resourceUsage := &types.NodeResourceArgs{}

	if err := json.Unmarshal([]byte(c.String("capacity")), resourceCapacity); err != nil {
		logrus.Errorf("[cmdAlloc] invalid resource capacity, err: %v", err)
		return err
	}
	if err := json.Unmarshal([]byte(c.String("usage")), resourceUsage); err != nil {
		logrus.Errorf("[cmdAlloc] invalid resource usage, err: %v", err)
		return err
	}

	err = volume.SetNodeResourceInfo(c.Context, node, resourceCapacity, resourceUsage)
	if err != nil {
		logrus.Errorf("[cmdSetNodeResourceInfo] failed to set node resource info, err: %v", err)
		return err
	}
	return nil
}

func cmdGetNode(c *cli.Context) error {
	// volume get-node --node xxx-node --workload-map {"workload-1": {"cpu": 100}} --fix
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	fix := c.Bool("fix")
	workloadMap := map[string]*types.WorkloadResourceArgs{}

	if err := json.Unmarshal([]byte(c.String("workload-map")), &workloadMap); err != nil {
		logrus.Errorf("[cmdGetNodeResourceInfo] invalid workload map, err: %v", err)
		return err
	}

	nodeResourceInfo, diffs, err := volume.GetNodeResourceInfo(c.Context, node, workloadMap, fix)
	if err != nil {
		logrus.Errorf("[cmdGetNodeResourceInfo] failed to get node resource info, err: %v", err)
		return err
	}

	printResult(map[string]interface{}{
		"resource_info": nodeResourceInfo,
		"diffs": diffs,
	})
	return nil
}