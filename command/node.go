package command

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateAddNodeCommand(pluginFactory func(*cli.Context) (interface{}, error), nodeResourceOpts Parsable) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()
		node := c.String("node")

		if err := nodeResourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
			logrus.Errorf("[cmdAddNode] invalid resource opts, err: %v", err)
			return err
		}

		out, err := call(plugin, "AddNode", c.Context, node, nodeResourceOpts)
		if err != nil {
			logrus.Errorf("[doAddNode] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[1] != nil {
			err = out[1].(error)
			logrus.Errorf("[doAddNode] plugin %v failed to add node, err: %v", pluginName, err)
			return err
		}

		printResult(out[0])
		return nil
	}

	return &cli.Command{
		Name:   "add-node",
		Usage:  "add node",
		Action: action,
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
}

func generateRemoveNodeCommand(pluginFactory func(*cli.Context) (interface{}, error)) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()
		node := c.String("node")

		out, err := call(plugin, "RemoveNode", c.Context, node)
		if err != nil {
			logrus.Errorf("[doRemoveNode] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[0] != nil {
			err = out[0].(error)
			logrus.Errorf("[doRemoveNode] plugin %v failed to remove node, err: %v", pluginName, err)
			return err
		}
		return nil
	}

	return &cli.Command{
		Name:   "remove-node",
		Usage:  "remove node",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "node name",
				Required: true,
			},
		},
	}
}
