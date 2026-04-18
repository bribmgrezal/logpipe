package dedup

import (
	"encoding/json"
	"os"
	"time"
)

// Config holds deduplication settings.
type Config struct {
	TTLSeconds int    `json:"ttl_seconds"`
	Field      string `json:"field"`
}

// LoadConfig reads a JSON config file for the dedup module.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// NewFromConfig constructs a Deduplicator from a Config.
func NewFromConfig(cfg *Config) *Deduplicator {
	ttl := time.Duration(cfg.TTLSeconds) * time.Second
	return New(ttl, cfg.Field)
}
