package tail

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds configuration for the Tailer.
type Config struct {
	Path         string `json:"path"`
	PollMs       int    `json:"poll_ms"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("tail: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("tail: parse config: %w", err)
	}
	if cfg.Path == "" {
		return nil, fmt.Errorf("tail: config missing 'path'")
	}
	return &cfg, nil
}

// NewFromConfig creates a Tailer from a Config.
func NewFromConfig(cfg *Config) *Tailer {
	poll := time.Duration(cfg.PollMs) * time.Millisecond
	return New(cfg.Path, poll)
}
