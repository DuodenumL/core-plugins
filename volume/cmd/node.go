package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/volume/types"
)

var addNodeCommand = &cli.Command{
	Name:   "add-node",
	Usage:  "add node",
	Action: cmdAddNode,
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
	},
}

var removeNodeCommand = &cli.Command{
	Name:   "remove-node",
	Usage:  "remove node",
	Action: cmdRemoveNode,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "node",
			Usage:    "node name",
			Required: true,
		},
	},
}

func cmdAddNode(c *cli.Context) error {
	// volume add-node --node xxx-node --resource-opts {"cpu": 1}
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	resourceOpts := &types.NodeResourceOpts{}

	if err := resourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
		logrus.Errorf("[cmdAddNode] invalid resource opts, err: %v", err)
		return err
	}

	resourceInfo, err := volume.AddNode(c.Context, node, resourceOpts)
	if err != nil {
		logrus.Errorf("[cmdAddNode] failed to add node, err: %v", err)
		return err
	}

	printResult(map[string]interface{}{
		"capacity": resourceInfo.Capacity,
		"usage":    resourceInfo.Usage,
	})

	return nil
}

func cmdRemoveNode(c *cli.Context) error {
	// volume add-node --node xxx-node
	volume, err := newVolume(c)
	if err != nil {
		return err
	}

	node := c.String("node")
	if err = volume.RemoveNode(c.Context, node); err != nil {
		logrus.Errorf("[cmdRemoveNode] failed to remove node, err: %v", err)
		return err
	}
	return nil
}
