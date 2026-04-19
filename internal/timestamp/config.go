package timestamp

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds configuration for the timestamp processor.
type Config struct {
	Field  string `json:"field"`
	InFmt  string `json:"in_format"`
	OutFmt string `json:"out_format"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("timestamp config: read file: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("timestamp config: parse: %w", err)
	}
	if cfg.Field == "" {
		return nil, fmt.Errorf("timestamp config: field is required")
	}
	return &cfg, nil
}

// NewFromConfig constructs a Processor from a Config.
func NewFromConfig(cfg *Config) (*Processor, error) {
	if cfg == nil {
		return nil, fmt.Errorf("timestamp config: nil config")
	}
	return New(cfg.Field, cfg.InFmt, cfg.OutFmt)
}
