package aggregate

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config holds aggregation settings loaded from a JSON file.
type Config struct {
	Field    string `json:"field"`
	Interval string `json:"interval"`
}

// LoadConfig reads an aggregation config from path.
func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}
	if c.Field == "" {
		return nil, fmt.Errorf("aggregate: field is required")
	}
	return &c, nil
}

// NewFromConfig builds an Aggregator from a Config.
func NewFromConfig(c *Config, out func([]byte) error) (*Aggregator, error) {
	interval := 10 * time.Second
	if c.Interval != "" {
		d, err := time.ParseDuration(c.Interval)
		if err != nil {
			return nil, fmt.Errorf("aggregate: invalid interval %q: %w", c.Interval, err)
		}
		interval = d
	}
	return New(c.Field, interval, out), nil
}
