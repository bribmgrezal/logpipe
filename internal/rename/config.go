package rename

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds the configuration for the rename module.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads and parses a rename config file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("rename: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("rename: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Renamer from a Config.
// Returns a no-op Renamer if cfg is nil.
func NewFromConfig(cfg *Config) *Renamer {
	if cfg == nil {
		return New(nil)
	}
	return New(cfg.Rules)
}
