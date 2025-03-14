package config

import (
	"fmt"
	"os"
)

// HelmValues represents the structure of values.yaml

type Config struct {
	File string
}

func NewConfig(file string) *Config {
	return &Config{
		File: file,
	}
}

func (c *Config) ReadHelmValues() ([]byte, error) {
	// Read values.yaml
	valuesFile, err := os.ReadFile(c.File)
	if err != nil {
		return nil, fmt.Errorf("error reading %v", err)
	}

	return valuesFile, nil
}
