package cast

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the configuration for the cast module.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("cast: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("cast: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Caster from a Config.
func NewFromConfig(cfg *Config) *Caster {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
