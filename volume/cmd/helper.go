package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/projecteru2/core-plugins/volume/models"
	"github.com/projecteru2/core-plugins/volume/types"
)

func newVolume(c *cli.Context) (*models.Volume, error) {
	configPath := c.String("config")
	config := &types.Config{}
	if err := configor.Load(config, configPath); err != nil {
		logrus.Errorf("[newVolume] failed to load config, err: %v", err)
		return nil, err
	}

	cpuMem, err := models.NewVolume(config)
	if err != nil {
		logrus.Errorf("[newVolume] failed to init volume, err: %v", err)
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
