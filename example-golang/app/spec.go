package app

import (
	"fmt"
	"os"

	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type MetaConfig struct {
	Lang  string       `yaml:"lang"`
	Code  string       `yaml:"code"`
	Label string       `yaml:"label"`
	Tag   language.Tag `yaml:"-"`
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

	t, err := language.Parse(config.Lang)
	if err != nil {
		return nil, fmt.Errorf(`spec: invalid language code "%s": %w (%s)`, config.Code, err, filepath)
	}
	config.Tag = t
	return &config, nil
}
