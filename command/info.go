package command

import (
	"encoding/json"
	"reflect"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func generateGetNodeResourceInfoCommand(pluginFactory func(*cli.Context) (interface{}, error), workloadMap Parsable) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()
		node := c.String("node")
		fix := c.Bool("fix")

		if err := workloadMap.ParseFromString(c.String("workload-map")); err != nil {
			logrus.Errorf("[cmdGetNodeResourceInfo] invalid workload map, err: %v", err)
			return err
		}

		out, err := call(plugin, "GetNodeResourceInfo", c.Context, node, workloadMap, fix)
		if err != nil {
			logrus.Errorf("[doGetReallocArgs] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doGetReallocArgs] plugin %v failed to get realloc args, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"resource_info": out[0],
			"diffs":         out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "get-node",
		Usage:  "get the total amount and usage of node resources",
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
				Value:    "{}",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "fix",
				Usage:    "fix the resource usage of node according to the resource usage of workload",
				Required: false,
				Value:    false,
			},
		},
	}
}

func generateSetNodeResourceInfoCommand(pluginFactory func(*cli.Context) (interface{}, error), resourceCapacity interface{}, resourceUsage interface{}) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()
		node := c.String("node")

		if err := json.Unmarshal([]byte(c.String("capacity")), resourceCapacity); err != nil {
			logrus.Errorf("[doSetReallocArgs] invalid resource capacity, err: %v", err)
			return err
		}
		if err := json.Unmarshal([]byte(c.String("usage")), resourceUsage); err != nil {
			logrus.Errorf("[doSetReallocArgs] invalid resource usage, err: %v", err)
			return err
		}

		out, err := call(plugin, "SetNodeResourceInfo", c.Context, node, resourceCapacity, resourceUsage)
		if err != nil {
			logrus.Errorf("[doSetReallocArgs] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doSetReallocArgs] plugin %v failed to set realloc args, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"resource_info": out[0],
			"diffs":         out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "set-node",
		Usage:  "set the total amount and usage of node resources",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "node name",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "capacity",
				Usage:    "resource capacity",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "usage",
				Usage:    "resource usage",
				Required: true,
			},
		},
	}
}

func generateSetNodeResourceUsageCommand(pluginFactory func(*cli.Context) (interface{}, error), nodeResourceOpts Parsable, nodeResourceArgs interface{}, workloadResourceArgs interface{}) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()
		node := c.String("node")
		delta := c.Bool("delta")
		incr := !c.Bool("decr")

		var out []interface{}

		if c.IsSet("node-resource-opts") {
			if err = nodeResourceOpts.ParseFromString(c.String("node-resource-opts")); err != nil {
				logrus.Errorf("[doSetNodeResourceUsage] plugin %v failed to parse node resource opts, err: %v", pluginName, err)
				return err
			}
			out, err = call(plugin, "SetNodeResourceUsage", c.Context, node, nodeResourceOpts, reflect.Zero(reflect.TypeOf(nodeResourceArgs)), reflect.Zero(makeSliceOfType(workloadResourceArgs, 0).Type()), delta, incr)

		} else if c.IsSet("node-resource-args") {
			if err = json.Unmarshal([]byte(c.String("node-resource-args")), nodeResourceArgs); err != nil {
				logrus.Errorf("[doSetNodeResourceUsage] plugin %v failed to parse node resource args, err: %v", pluginName, err)
				return err
			}
			out, err = call(plugin, "SetNodeResourceUsage", c.Context, node, makeTypedNil(nodeResourceOpts), nodeResourceArgs, makeTypedNil(makeSliceOfType(workloadResourceArgs, 0).Interface()), delta, incr)

		} else if c.IsSet("workload-resource-args") {
			workloadResourceArgsSlice := makeSliceOfType(workloadResourceArgs, len(c.StringSlice("workload-resource-args")))
			for _, args := range c.StringSlice("workload-resource-args") {
				resourceArgs := newPtrOfType(workloadResourceArgs)
				if err := json.Unmarshal([]byte(args), resourceArgs); err != nil {
					logrus.Errorf("[doSetNodeResourceUsage] invalid workload resource args, err: %v", err)
					return err
				}
				workloadResourceArgsSlice = reflect.Append(workloadResourceArgsSlice, reflect.ValueOf(resourceArgs))
			}
			out, err = call(plugin, "SetNodeResourceUsage", c.Context, node, makeTypedNil(nodeResourceOpts), makeTypedNil(nodeResourceArgs), workloadResourceArgsSlice.Interface(), delta, incr)

		} else {
			logrus.Errorf("[doSetNodeResourceUsage] plugin %v receives invalid parameters", pluginName)
			return ErrInvalidParams
		}

		if err != nil {
			logrus.Errorf("[doSetNodeResourceUsage] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doSetNodeResourceUsage] plugin %v failed to set node resource usage, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"before": out[0],
			"after":  out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "set-node-usage",
		Usage:  "set the usage of node resources",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "node name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "node-resource-opts",
				Usage: "node resource opts",
			},
			&cli.StringFlag{
				Name:  "node-resource-args",
				Usage: "node resource args",
			},
			&cli.StringSliceFlag{
				Name:  "workload-resource-args",
				Usage: "workload resource args",
			},
			&cli.BoolFlag{
				Name:  "delta",
				Usage: "to use delta changes instead of absolute changes",
			},
			&cli.BoolFlag{
				Name:  "decr",
				Usage: "to decrease the resource usage",
			},
		},
	}
}

func generateSetNodeResourceCapacityCommand(pluginFactory func(*cli.Context) (interface{}, error), nodeResourceOpts Parsable, nodeResourceArgs interface{}) *cli.Command {
	action := func(c *cli.Context) error {
		plugin, err := pluginFactory(c)
		if err != nil {
			logrus.Errorf("failed to init plugin, err: %v", err)
			return err
		}
		pluginName := reflect.Indirect(reflect.ValueOf(plugin)).Type().Name()
		node := c.String("node")
		delta := c.Bool("delta")
		incr := !c.Bool("decr")

		var out []interface{}

		if c.IsSet("node-resource-opts") {
			if err = nodeResourceOpts.ParseFromString(c.String("node-resource-opts")); err != nil {
				logrus.Errorf("[doSetNodeResourceCapacity] plugin %v failed to parse node resource opts, err: %v", pluginName, err)
				return err
			}
			out, err = call(plugin, "SetNodeResourceCapacity", c.Context, node, nodeResourceOpts, makeTypedNil(nodeResourceArgs), delta, incr)

		} else if c.IsSet("node-resource-args") {
			if err = json.Unmarshal([]byte(c.String("node-resource-args")), nodeResourceArgs); err != nil {
				logrus.Errorf("[doSetNodeResourceCapacity] plugin %v failed to parse node resource args, err: %v", pluginName, err)
				return err
			}
			out, err = call(plugin, "SetNodeResourceCapacity", c.Context, node, makeTypedNil(nodeResourceOpts), nodeResourceArgs, delta, incr)

		} else {
			logrus.Errorf("[doSetNodeResourceCapacity] plugin %v receives invalid parameters", pluginName)
			return ErrInvalidParams
		}

		if err != nil {
			logrus.Errorf("[doSetNodeResourceCapacity] failed to call plugin %v, err: %v", pluginName, err)
			return err
		}
		if out[2] != nil {
			err = out[2].(error)
			logrus.Errorf("[doSetNodeResourceCapacity] plugin %v failed to set node resource capacity, err: %v", pluginName, err)
			return err
		}

		printResult(map[string]interface{}{
			"before": out[0],
			"after":  out[1],
		})
		return nil
	}

	return &cli.Command{
		Name:   "set-node-capacity",
		Usage:  "set the total amount of node resources",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "node",
				Usage:    "node name",
				Required: true,
			},
			&cli.StringFlag{
				Name:  "node-resource-opts",
				Usage: "node resource opts",
			},
			&cli.StringFlag{
				Name:  "node-resource-args",
				Usage: "node resource args",
			},
			&cli.BoolFlag{
				Name:  "delta",
				Usage: "to use delta changes instead of absolute changes",
			},
			&cli.BoolFlag{
				Name:  "decr",
				Usage: "to decrease the resource usage",
			},
		},
	}
}
