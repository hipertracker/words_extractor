package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type MetaConfig struct {
	Lang  string `yaml:"lang"`
	Code  string `yaml:"code"`
	Label string `yaml:"label"`
}

func ReadSpec(filepath string) (*MetaConfig, error) {
	p, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf(`spec: reading YAML file "%s": %w`, filepath, err)
	}

	var config MetaConfig
	if err = yaml.Unmarshal(p, &config); err != nil {
		return nil, fmt.Errorf(`spec: parsing YAML file "%s": %w`, filepath, err)
	}

	return &config, nil
}
