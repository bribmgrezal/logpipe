package router

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds router configuration loaded from a JSON file.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("router config: open %q: %w", path, err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("router config: decode: %w", err)
	}
	return &cfg, nil
}

// Validate checks that all rules have required fields set.
func (c *Config) Validate() error {
	for i, r := range c.Rules {
		if r.Field == "" {
			return fmt.Errorf("router config: rule[%d] missing field", i)
		}
		if r.Target == "" {
			return fmt.Errorf("router config: rule[%d] missing target", i)
		}
	}
	return nil
}
