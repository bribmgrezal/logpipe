package checkpoint

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds configuration for the checkpoint module.
type Config struct {
	Path string `json:"path"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("checkpoint: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("checkpoint: parse config: %w", err)
	}
	if cfg.Path == "" {
		return nil, fmt.Errorf("checkpoint: config missing required field 'path'")
	}
	return &cfg, nil
}

// NewFromConfig constructs a Checkpoint from a Config.
func NewFromConfig(cfg *Config) (*Checkpoint, error) {
	if cfg == nil {
		return nil, fmt.Errorf("checkpoint: nil config")
	}
	return New(cfg.Path)
}
