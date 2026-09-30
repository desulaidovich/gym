package config

import (
	"fmt"
	"os"

	"go.uber.org/config"
)

func Parse(configPath string) (*Config, error) {
	provider, err := config.NewYAML(
		config.File(configPath),
		config.Expand(os.LookupEnv),
	)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	var cfg Config
	if err := provider.Get(config.Root).Populate(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return &cfg, nil
}
