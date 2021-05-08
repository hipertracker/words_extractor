package main

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type metaConfig struct {
	Lang   string
	Code   string
	Label  string
	Year   string
	Strong bool
}

func (m *metaConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, m)
}

func getMeta(path string) metaConfig {
	data, err := ioutil.ReadFile(path)
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
