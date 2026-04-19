package batch

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds batcher configuration.
type Config struct {
	Size     int    `json:"size"`
	Interval string `json:"interval"`
}

// LoadConfig reads a JSON config file for the batcher.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("batch: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("batch: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig constructs a Batcher from a Config.
func NewFromConfig(cfg *Config, flushFn func([]map[string]any)) (*Batcher, error) {
	if cfg == nil {
		return nil, fmt.Errorf("batch: nil config")
	}
	interval := 5 * time.Second
	if cfg.Interval != "" {
		d, err := time.ParseDuration(cfg.Interval)
		if err != nil {
			return nil, fmt.Errorf("batch: invalid interval: %w", err)
		}
		interval = d
	}
	return New(cfg.Size, interval, flushFn), nil
}
