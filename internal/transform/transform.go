package transform

import (
	"encoding/json"
	"fmt"
)

// Rule defines a single transformation to apply to a log record.
type Rule struct {
	Action string // "rename", "delete", "add"
	Field  string
	Value  string // used by "rename" (new name) and "add" (value)
}

// Transformer applies a set of rules to JSON log lines.
type Transformer struct {
	rules []Rule
}

// New creates a new Transformer with the given rules.
func New(rules []Rule) *Transformer {
	return &Transformer{rules: rules}
}

// Apply transforms a raw JSON log line according to the configured rules.
// Returns the transformed JSON or an error if the input is not valid JSON.
func (t *Transformer) Apply(line string) (string, error) {
	if len(t.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return "", fmt.Errorf("transform: invalid JSON: %w", err)
	}

	for _, rule := range t.rules {
		switch rule.Action {
		case "rename":
			if val, ok := record[rule.Field]; ok {
				record[rule.Value] = val
				delete(record, rule.Field)
			}
		case "delete":
			delete(record, rule.Field)
		case "add":
			record[rule.Field] = rule.Value
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("transform: marshal error: %w", err)
	}
	return string(out), nil
}
