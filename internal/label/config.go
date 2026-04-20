package label

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds label configuration loaded from a file.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("label: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("label: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Labeler from a Config.
func NewFromConfig(cfg *Config) *Labeler {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
