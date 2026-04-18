package enrich

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds enrichment configuration.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads an enrichment config from a JSON file.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("enrich: read config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("enrich: parse config: %w", err)
	}
	return &cfg, nil
}

// NewFromConfig creates an Enricher from a Config.
func NewFromConfig(cfg *Config) *Enricher {
	return New(cfg.Rules)
}
