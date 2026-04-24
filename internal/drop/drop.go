package drop

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a condition for dropping a log line.
type Rule struct {
	Field    string `json:"field"`
	Operator string `json:"operator"` // eq, contains, exists, missing
	Value    string `json:"value"`
}

// Dropper discards log lines matching any configured rule.
type Dropper struct {
	rules []Rule
}

// New creates a Dropper with the given rules.
func New(rules []Rule) *Dropper {
	return &Dropper{rules: rules}
}

// Apply returns an empty string if the line matches any drop rule,
// otherwise returns the original line unchanged.
func (d *Dropper) Apply(line string) (string, error) {
	if len(d.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return "", fmt.Errorf("drop: invalid JSON: %w", err)
	}

	for _, rule := range d.rules {
		if matchRule(record, rule) {
			return "", nil
		}
	}
	return line, nil
}

func matchRule(record map[string]interface{}, rule Rule) bool {
	val, exists := record[rule.Field]
	switch rule.Operator {
	case "exists":
		return exists
	case "missing":
		return !exists
	case "eq":
		return exists && fmt.Sprintf("%v", val) == rule.Value
	case "contains":
		return exists && strings.Contains(fmt.Sprintf("%v", val), rule.Value)
	}
	return false
}
