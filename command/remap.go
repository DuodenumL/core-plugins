package command

import (
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateGetRemapArgsCommand(pluginFactory func(*cli.Context) (interface{}, error), workloadMap Parsable) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()

		node := c.String("node")

		if err := workloadMap.ParseFromString(c.String("workload-map")); err != nil {
			logrus.Errorf("[doGetRemapArgs] invalid workload map, err: %v", err)
			return err
		}

		out, err := call(plugin, "GetRemapArgs", c.Context, node, workloadMap)
		if err != nil {
			logrus.Errorf("[doGetRemapArgs] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[1] != nil {
			err = out[1].(error)
			logrus.Errorf("[doGetRemapArgs] plugin %v failed to get remap args, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"engine_args_map": out[0],
		})
		return nil
	}

	return &cli.Command{
		Name:   "get-remap-args",
		Usage:  "remap (pure calculation)",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "node name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "workload-map",
				Usage:    "the resource usage of all workloads belonging to this node",
				Required: true,
			},
		},
	}
}
