package dedupe

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds configuration for the Deduper.
type Config struct {
	Fields  []string `json:"fields"`
	MaxSize int      `json:"max_size"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("dedupe: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("dedupe: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Deduper from a Config.
func NewFromConfig(cfg *Config) *Deduper {
	if cfg == nil {
		return New(nil, 0)
	}
	return New(cfg.Fields, cfg.MaxSize)
}
