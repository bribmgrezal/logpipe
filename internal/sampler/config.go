package sampler

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds sampler configuration loaded from a file.
type Config struct {
	Rate uint64 `json:"rate"`
}

// LoadConfig reads a JSON config file and returns a Config.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("sampler: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("sampler: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates a Sampler from a Config.
func NewFromConfig(cfg *Config) *Sampler {
	return New(cfg.Rate)
}
