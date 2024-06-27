package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	RUN_MODE string `yaml:"RUN_MODE"`
}

func LoadConfig(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(buf, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
