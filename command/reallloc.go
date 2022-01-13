package command

import (
	"encoding/json"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateGetReallocArgsCommand(pluginFactory func(*cli.Context) (interface{}, error), oldResourceArgs interface{}, newResourceOpts Parsable) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()

		node := c.String("node")

		if err := json.Unmarshal([]byte(c.String("old")), oldResourceArgs); err != nil {
			logrus.Errorf("[doGetReallocArgs] invalid old resource args, err: %v", err)
			return err
		}
		if err := newResourceOpts.ParseFromString(c.String("resource-opts")); err != nil {
			logrus.Errorf("[doGetReallocArgs] invalid new resource opts, err: %v", err)
			return err
		}

		out, err := call(plugin, "GetReallocArgs", c.Context, node, oldResourceArgs, newResourceOpts)
		if err != nil {
			logrus.Errorf("[doGetReallocArgs] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[3] != nil {
			err = out[3].(error)
			logrus.Errorf("[doGetReallocArgs] plugin %v failed to get realloc args, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"engine_args":   out[0],
			"delta":         out[1],
			"resource_args": out[2],
		})
		return nil
	}

	return &cli.Command{
		Name:   "get-realloc-args",
		Usage:  "reallocate resources (pure calculation)",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "node name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "old",
				Usage:    "origin resource usage",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "resource-opts",
				Usage:    "resource options in delta form",
				Required: true,
			},
		},
	}
}
