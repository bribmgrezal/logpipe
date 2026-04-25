package limit

import (
	"encoding/json"
	"fmt"
)

// Rule defines a field and its maximum allowed byte length.
type Rule struct {
	Field  string `json:"field"`
	MaxLen int    `json:"max_len"`
	Drop   bool   `json:"drop"` // if true, drop the entire line instead of truncating
}

// Limiter enforces byte-length limits on JSON log fields.
type Limiter struct {
	rules []Rule
}

// New creates a Limiter with the given rules.
func New(rules []Rule) *Limiter {
	return &Limiter{rules: rules}
}

// Apply enforces field length limits on a JSON log line.
// Returns the (possibly modified) line, a drop flag, and any error.
func (l *Limiter) Apply(line []byte) ([]byte, bool, error) {
	if len(l.rules) == 0 {
		return line, false, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, false, fmt.Errorf("limit: invalid JSON: %w", err)
	}

	for _, rule := range l.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		str, ok := val.(string)
		if !ok {
			continue
		}
		if len(str) > rule.MaxLen {
			if rule.Drop {
				return nil, true, nil
			}
			record[rule.Field] = str[:rule.MaxLen]
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, false, fmt.Errorf("limit: marshal error: %w", err)
	}
	return out, false, nil
}
