package config

import (
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Database struct {
		Url    string `yaml:"url"`
		Driver string `yaml:"driver"`
	} `yaml:"database"`
	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Id   string `yaml:"id"`
	} `yaml:"server"`
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	config := Config{}
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
