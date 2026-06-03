package main

import (
	"log"
	"os"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Logger  LoggerConf  `yaml:"logger"`
	Listen  ListenConf  `yaml:"listen"`
	Storage StorageConf `yaml:"storage"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}

type ListenConf struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	GRPCPort string `yaml:"grpc_port"`
}

type StorageConf struct {
	Type  string `yaml:"type"`
	DBURL string `yaml:"dbUrl"`
}

func NewConfig(path string) Config {
	var config Config

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}

	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		f.Close()
		log.Fatalf("failed to decode config: %v", err)
	}

	f.Close()
	return config
}
