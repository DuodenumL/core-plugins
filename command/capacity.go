package command

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateGetNodesDeployCapacityCommand(pluginFactory func(*cli.Context) (interface{}, error), resourceOpts Parsable) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()

		nodes := c.StringSlice("node")

		if err := resourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
			logrus.Errorf("[cmdCapacity] invalid resource opts, err: %v", err)
			return err
		}

		out, err := call(plugin, "GetNodesDeployCapacity", c.Context, nodes, resourceOpts)
		if err != nil {
			logrus.Errorf("[doGetNodesDeployCapacity] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doGetNodesDeployCapacity] plugin %v failed to get deploy capacity, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"nodes": out[0],
			"total": out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "get-capacity",
		Usage:  "calculate how many workloads can be allocated according to the specified resource options (pure calculation)",
		Action: action,
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
}
