package flatten

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds configuration for the Flattener.
type Config struct {
	Separator string `json:"separator"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("flatten: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("flatten: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig returns a Flattener from a Config.
func NewFromConfig(cfg *Config) *Flattener {
	if cfg == nil {
		return New(".")
	}
	return New(cfg.Separator)
}
