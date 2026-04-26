package tag

import (
	"encoding/json"
	"fmt"
)

// Rule defines a tagging rule: if field matches value, add tag to the tags field.
type Rule struct {
	Field  string `json:"field"`
	Match  string `json:"match"`
	Tag    string `json:"tag"`
	Target string `json:"target"` // defaults to "tags"
}

// Tagger applies tagging rules to log lines.
type Tagger struct {
	rules []Rule
}

// New creates a new Tagger with the given rules.
func New(rules []Rule) *Tagger {
	return &Tagger{rules: rules}
}

// Apply processes a JSON log line, appending tags to the target field.
func (t *Tagger) Apply(line string) (string, error) {
	if len(t.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("tag: invalid JSON: %w", err)
	}

	for _, r := range t.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", val) != r.Match {
			continue
		}
		target := r.Target
		if target == "" {
			target = "tags"
		}
		switch existing := obj[target].(type) {
		case []interface{}:
			obj[target] = append(existing, r.Tag)
		case nil:
			obj[target] = []interface{}{r.Tag}
		default:
			obj[target] = []interface{}{existing, r.Tag}
		}
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("tag: marshal error: %w", err)
	}
	return string(out), nil
}
