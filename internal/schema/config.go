package schema

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds schema validation configuration.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a schema config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("schema: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("schema: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Validator from a Config.
func NewFromConfig(cfg *Config) *Validator {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
