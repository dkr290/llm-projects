package config

import (
	"fmt"
	"os"
)

// HelmValues represents the structure of values.yaml

type Config struct {
	ValuesFile string
	ChartFile  string
}

func NewConfig(values, chart string) *Config {
	return &Config{
		ValuesFile: values,
		ChartFile:  chart,
	}
}

func (c *Config) ReadHelmValues() ([]byte, []byte, error) {
	// Read values.yaml
	valuesFile, err := os.ReadFile(c.ValuesFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading %v", err)
	}

	chartFile, err := os.ReadFile(c.ChartFile)
	if err != nil {
		return nil, nil, fmt.Errorf("error reading %v", err)
	}

	return valuesFile, chartFile, nil
}
