package config

import (
	"Skland/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Account models.AccountInfo `yaml:"account"`
}

func LoadConfig() (*Config, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}
	fmt.Printf("exePath: %v\n", exePath)

	configPath := filepath.Join(filepath.Dir(exePath), "config.yaml")
	fmt.Printf("configPath: %v\n", configPath)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if cfg.Account.Phone == "" || cfg.Account.Password == "" {
		return nil, fmt.Errorf("missing required account configuration")
	}

	return &cfg, nil
}
