package command

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateGetMostIdleNodeCommand(pluginFactory func(*cli.Context) (interface{}, error)) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()

		nodes := c.StringSlice("node")

		out, err := call(plugin, "GetMostIdleNode", c.Context, nodes)
		if err != nil {
			logrus.Errorf("[doGetMostIdleNode] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doGetMostIdleNode] plugin %v failed to get remap args, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"node":     out[0],
			"priority": out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "get-idle",
		Usage:  "get most idle nodes",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "node",
				Usage:    "node names",
				Required: true,
			},
		},
	}
}
