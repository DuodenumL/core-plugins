package cmd

import "github.com/urfave/cli/v2"

var Commands = []*cli.Command{
	allocCommand,
	capacityCommand,
	getNodeCommand,
	setNodeCommand,
	addNodeCommand,
	removeNodeCommand,
	reallocCommand,
	remapCommand,
	addNodeCommand,
	removeNodeCommand,
	updateCapacityCommand,
	updateUsageCommand,
}
