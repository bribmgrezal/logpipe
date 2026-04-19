package normalize

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a normalization rule for a single field.
type Rule struct {
	Field string `json:"field"`
	Op    string `json:"op"` // lowercase, uppercase, trim
}

// Normalizer applies normalization rules to JSON log lines.
type Normalizer struct {
	rules []Rule
}

// New creates a Normalizer with the given rules.
func New(rules []Rule) *Normalizer {
	return &Normalizer{rules: rules}
}

// Apply normalizes fields in a JSON log line and returns the result.
func (n *Normalizer) Apply(line []byte) ([]byte, error) {
	if len(n.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("normalize: invalid JSON: %w", err)
	}

	for _, rule := range n.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		str, ok := val.(string)
		if !ok {
			continue
		}
		switch rule.Op {
		case "lowercase":
			record[rule.Field] = strings.ToLower(str)
		case "uppercase":
			record[rule.Field] = strings.ToUpper(str)
		case "trim":
			record[rule.Field] = strings.TrimSpace(str)
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("normalize: marshal error: %w", err)
	}
	return out, nil
}
