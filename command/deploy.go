package command

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateGetDeployArgsCommand(pluginFactory func(*cli.Context) (interface{}, error), resourceOpts Parsable) *cli.Command {
	action := func(c *cli.Context) error {
		node := c.String("node")
		deployCount := c.Int("deploy")
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()

		if err := resourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
			logrus.Errorf("[doGetDeployArgs] plugin %v invalid resource opts, err: %v", pluginName, err)
			return err
		}

		out, err := call(plugin, "GetDeployArgs", c.Context, node, deployCount, resourceOpts)
		if err != nil {
			logrus.Errorf("[doGetDeployArgs] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doGetDeployArgs] plugin %v failed to get deploy args, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"engine_args":   out[0],
			"resource_args": out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "get-deploy-args",
		Usage:  "get deploy args",
		Action: action,
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
}
