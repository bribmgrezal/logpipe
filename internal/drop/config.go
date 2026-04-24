package drop

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the configuration for the drop module.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads and parses a JSON config file for the drop module.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("drop: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("drop: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Dropper from a Config.
// If cfg is nil, a no-op Dropper is returned.
func NewFromConfig(cfg *Config) *Dropper {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
