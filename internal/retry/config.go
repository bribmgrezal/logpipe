package retry

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds retry configuration loaded from a JSON file.
type Config struct {
	MaxAttempts int `json:"max_attempts"`
	DelayMs     int `json:"delay_ms"`
}

// Validate checks that the Config fields have sensible values.
// It returns an error if MaxAttempts is less than 1 or DelayMs is negative.
func (c *Config) Validate() error {
	if c.MaxAttempts < 1 {
		return fmt.Errorf("retry: max_attempts must be at least 1, got %d", c.MaxAttempts)
	}
	if c.DelayMs < 0 {
		return fmt.Errorf("retry: delay_ms must be non-negative, got %d", c.DelayMs)
	}
	return nil
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
	if err := cfg.Validate(); err != nil {
		return nil, err
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
