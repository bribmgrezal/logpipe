package redact

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config holds redaction configuration loaded from a file.
type Config struct {
	Rules []Rule `json:"rules"`
}

// LoadConfig reads a JSON config file and returns a Redactor.
func LoadConfig(path string) (*Redactor, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("redact: open config: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("redact: decode config: %w", err)
	}

	if len(cfg.Rules) == 0 {
		return nil, fmt.Errorf("redact: config %q contains no rules", path)
	}

	return New(cfg.Rules), nil
}
