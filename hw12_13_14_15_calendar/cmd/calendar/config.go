package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config конфигурация приложения.
type Config struct {
	Logger  LoggerConfig  `yaml:"logger"`
	Storage StorageConfig `yaml:"storage"`
	Server  ServerConfig  `yaml:"server"`
}

// LoggerConfig конфигурация логгера.
type LoggerConfig struct {
	Level string `yaml:"level"` // Уровень логирования.
}

// StorageConfig конфигурация хранилища.
type StorageConfig struct {
	UseInMemory    bool   `yaml:"user_in_memory"`   // Использовать хранилище в оперативной памяти?
	DataSourceName string `yaml:"data_source_name"` // Строка подключения к БД постгрес.
}

// ServerConfig конфигурация для запуска http сервера.
type ServerConfig struct {
	Host string `yaml:"host"` // Адрес.
	Port int    `yaml:"port"` // Порт.
}

func NewConfigFomFile(path string) (cfg *Config, err error) {
	c := GetDefault()
	yamlContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(yamlContent, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetDefault() *Config {
	c := &Config{
		Logger:  LoggerConfig{Level: "info"},
		Storage: StorageConfig{UseInMemory: true},
		Server:  ServerConfig{Host: "0.0.0.0", Port: 8080},
	}
	return c
}
