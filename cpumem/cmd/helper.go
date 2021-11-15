package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/cpumem/models"
	"github.com/projecteru2/core-plugins/cpumem/types"
)

func newCPUMem(c *cli.Context) (*models.CPUMem, error) {
	configPath := c.String("config")
	config := &types.Config{}
	if err := configor.Load(config, configPath); err != nil {
		logrus.Errorf("[newCPUMem] failed to load config, err: %v", err)
		return nil, err
	}

	cpuMem, err := models.NewCPUMem(config)
	if err != nil {
		logrus.Errorf("[newCPUMem] failed to init cpumem, err: %v", err)
		return nil, err
	}
	return cpuMem, nil
}

func printResult(result map[string]interface{}) {
	body, err := json.Marshal(result)
	if err != nil {
		logrus.Errorf("[cmdAlloc] failed to marshal result %+v, err: %v", result, err)
		fmt.Println("invalid json")
	}

	fmt.Println(string(body))
}
