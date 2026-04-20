package coalesce

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the file-level configuration for the coalesce module.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("coalesce: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("coalesce: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Coalescer from a Config.
// Returns an error if cfg is nil.
func NewFromConfig(cfg *Config) (*Coalescer, error) {
	if cfg == nil {
		return nil, fmt.Errorf("coalesce: config must not be nil")
	}
	return New(cfg.Rules), nil
}
