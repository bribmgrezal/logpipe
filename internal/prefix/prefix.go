package prefix

import (
	"encoding/json"
	"fmt"
)

// Rule defines a field to prepend a string prefix to.
type Rule struct {
	Field  string `json:"field"`
	Prefix string `json:"prefix"`
}

// Prefixer applies prefix rules to JSON log lines.
type Prefixer struct {
	rules []Rule
}

// New creates a new Prefixer with the given rules.
func New(rules []Rule) *Prefixer {
	return &Prefixer{rules: rules}
}

// Apply prepends configured prefixes to specified fields in the JSON line.
// If no rules are configured, the line is returned unchanged.
// If the line is not valid JSON, an error is returned.
func (p *Prefixer) Apply(line string) (string, error) {
	if len(p.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return "", fmt.Errorf("prefix: invalid JSON: %w", err)
	}

	for _, rule := range p.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		str, ok := val.(string)
		if !ok {
			continue
		}
		record[rule.Field] = rule.Prefix + str
	}

	out, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("prefix: marshal error: %w", err)
	}
	return string(out), nil
}
