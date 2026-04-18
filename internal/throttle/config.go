package throttle

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds throttle configuration.
type Config struct {
	WindowSeconds int `json:"window_seconds"`
}

// LoadConfig reads a throttle config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("throttle: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("throttle: parse config: %w", err)
	}
	if cfg.WindowSeconds <= 0 {
		return nil, fmt.Errorf("throttle: window_seconds must be > 0")
	}
	return &cfg, nil
}

// NewFromConfig creates a Throttler from a Config.
func NewFromConfig(cfg *Config) *Throttler {
	if cfg == nil {
		return New(time.Second)
	}
	return New(time.Duration(cfg.WindowSeconds) * time.Second)
}
