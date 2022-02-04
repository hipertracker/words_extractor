package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type metaConfig struct {
	Lang  string
	Code  string
	Label string
}

func (m *metaConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, m)
}

func GetYAML(filepath string) metaConfig {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	config := metaConfig{}
	err = config.Parse(data)
	if err != nil {
		panic(err)
	}
	return config
}
