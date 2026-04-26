package tag

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the configuration for the tag module.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tag: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("tag: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Tagger from a Config.
func NewFromConfig(cfg *Config) *Tagger {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
