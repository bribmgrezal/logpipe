package coalesce

import (
	"encoding/json"
	"fmt"
)

// Rule defines a coalesce operation: pick the first non-empty value from
// Fields and write it into Target.
type Rule struct {
	Fields []string `json:"fields"`
	Target string   `json:"target"`
}

// Coalescer applies coalesce rules to JSON log lines.
type Coalescer struct {
	rules []Rule
}

// New creates a Coalescer with the given rules.
func New(rules []Rule) *Coalescer {
	return &Coalescer{rules: rules}
}

// Apply processes a single JSON log line, resolving coalesce rules.
// It returns the modified line or an error if the input is not valid JSON.
func (c *Coalescer) Apply(line []byte) ([]byte, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("coalesce: invalid JSON: %w", err)
	}

	for _, rule := range c.rules {
		if rule.Target == "" || len(rule.Fields) == 0 {
			continue
		}
		for _, field := range rule.Fields {
			val, ok := record[field]
			if !ok {
				continue
			}
			if s, ok := val.(string); ok && s == "" {
				continue
			}
			if val == nil {
				continue
			}
			record[rule.Target] = val
			break
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("coalesce: marshal error: %w", err)
	}
	return out, nil
}
