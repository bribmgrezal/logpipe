package retry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds retry configuration loaded from a JSON file.
type Config struct {
	MaxAttempts int    `json:"max_attempts"`
	DelayMs     int    `json:"delay_ms"`
}

// LoadConfig reads a retry config from the given JSON file path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("retry: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("retry: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Retryer from config and a writer function.
func NewFromConfig(cfg *Config, writer func([]byte) error) *Retryer {
	if cfg == nil {
		return New(1, 0, writer)
	}
	return New(cfg.MaxAttempts, time.Duration(cfg.DelayMs)*time.Millisecond, writer)
}
