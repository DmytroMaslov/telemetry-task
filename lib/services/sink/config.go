package sink

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	defaultBind = "localhost:8080"
	defaultFile = "./tmp/metrics.txt"
)

type Config struct {
	BindAddress   string `yaml:"bind_address"`
	FilePath      string `yaml:"file_path"`
	BufferSize    int    `yaml:"buffer_size"`
	FlushInterval int    `yaml:"flush_interval"`
	RateLimit     int    `yaml:"rate_limit"`
}

func DefaultConfig() *Config {
	return &Config{
		BindAddress:   defaultBind,
		FilePath:      defaultFile,
		BufferSize:    1024,
		FlushInterval: 100,
		RateLimit:     1024 * 1024,
	}
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("open config file: %v", err)
	}

	var cfg Config
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("unmarshal config file: %v", err)
	}
	return &cfg, nil
}
