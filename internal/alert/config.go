package alert

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds alert configuration loaded from a file.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads an alert config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("alert: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("alert: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates an Alerter from a Config and writer.
func NewFromConfig(cfg *Config, w interface{ Write([]byte) (int, error) }) *Alerter {
	if cfg == nil {
		return New(nil, w)
	}
	return New(cfg.Rules, w)
}
