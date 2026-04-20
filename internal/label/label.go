package label

import (
	"encoding/json"
	"fmt"
)

// Rule defines a static label to add under a given key.
type Rule struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Labeler attaches static key-value labels to each log line.
type Labeler struct {
	rules []Rule
}

// New creates a Labeler with the given rules.
func New(rules []Rule) *Labeler {
	return &Labeler{rules: rules}
}

// Apply adds configured labels to the JSON log line.
// Existing keys are not overwritten.
func (l *Labeler) Apply(line []byte) ([]byte, error) {
	if len(l.rules) == 0 {
		return line, nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("label: invalid JSON: %w", err)
	}
	for _, r := range l.rules {
		if _, exists := obj[r.Key]; !exists {
			obj[r.Key] = r.Value
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("label: marshal error: %w", err)
	}
	return out, nil
}
