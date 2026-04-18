package mask

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds masking configuration loaded from a file.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mask: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("mask: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Masker from a Config.
func NewFromConfig(cfg *Config) (*Masker, error) {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
