package format

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds formatter configuration.
type Config struct {
	Template string `json:"template"`
}

// LoadConfig reads a JSON config file for the formatter.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("format: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("format: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Formatter from a Config.
func NewFromConfig(cfg *Config) *Formatter {
	if cfg == nil {
		return New("")
	}
	return New(cfg.Template)
}
