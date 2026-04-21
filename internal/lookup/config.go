package lookup

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the configuration for the lookup module.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lookup: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("lookup: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Lookup from a Config.
func NewFromConfig(cfg *Config) (*Lookup, error) {
	if cfg == nil {
		return New(nil), nil
	}
	for i, r := range cfg.Rules {
		if r.Field == "" {
			return nil, fmt.Errorf("lookup: rule[%d] missing field", i)
		}
		if len(r.Table) == 0 {
			return nil, fmt.Errorf("lookup: rule[%d] has empty table", i)
		}
	}
	return New(cfg.Rules), nil
}
