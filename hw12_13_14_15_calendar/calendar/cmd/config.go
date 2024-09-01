package main

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Logger     LoggerConfig     `yaml:"logger"`
	Repository RepositoryConfig `yaml:"repository"`
	Server     ServerConfig     `yaml:"server"`
}

func (c *Config) check() error {
	if len(c.Logger.Level) == 0 {
		return errors.New("logger Level is required")
	}
	if !c.Repository.UseInMemory && c.Repository.DataSourceName == "" {
		return errors.New("data source name is required if in memory repository is not used")
	}
	if c.Server.Address == "" {
		return errors.New("server dddress is required")
	}
	return nil
}

type LoggerConfig struct {
	Level string `yaml:"level"`
}

type RepositoryConfig struct {
	UseInMemory    bool   `yaml:"useInMemory"`
	DataSourceName string `yaml:"dataSourceName"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
}

func readConfigFromYaml(filename string) (*Config, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err = yaml.Unmarshal(content, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
