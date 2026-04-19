package truncate

import (
	"encoding/json"
	"fmt"
)

// Rule defines a field and its max byte length.
type Rule struct {
	Field  string `json:"field"`
	MaxLen int    `json:"max_len"`
	Suffix string `json:"suffix"` // optional, default "..."
}

// Truncator trims string fields that exceed a maximum length.
type Truncator struct {
	rules []Rule
}

// New creates a Truncator with the given rules.
func New(rules []Rule) *Truncator {
	return &Truncator{rules: rules}
}

// Apply truncates configured fields in the JSON line and returns the result.
func (t *Truncator) Apply(line []byte) ([]byte, error) {
	if len(t.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return nil, fmt.Errorf("truncate: invalid json: %w", err)
	}

	for _, r := range t.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if len(s) > r.MaxLen {
			suffix := r.Suffix
			if suffix == "" {
				suffix = "..."
			}
			cut := r.MaxLen - len(suffix)
			if cut < 0 {
				cut = 0
			}
			obj[r.Field] = s[:cut] + suffix
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("truncate: marshal: %w", err)
	}
	return out, nil
}
