package merge

import (
	"encoding/json"
	"fmt"
)

// Rule defines a merge operation: take fields from `sources` and merge them
// into a single target field (a JSON object).
type Rule struct {
	Target  string   `json:"target"`
	Sources []string `json:"sources"`
	Remove  bool     `json:"remove"` // remove source fields after merge
}

// Merger merges multiple fields into a single nested object field.
type Merger struct {
	rules []Rule
}

// New creates a Merger with the given rules.
func New(rules []Rule) *Merger {
	return &Merger{rules: rules}
}

// Apply merges source fields into target fields per each rule.
// Returns the modified JSON line, or an error if the input is invalid JSON.
func (m *Merger) Apply(line []byte) ([]byte, error) {
	if len(m.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("merge: invalid JSON: %w", err)
	}

	for _, rule := range m.rules {
		if rule.Target == "" || len(rule.Sources) == 0 {
			continue
		}

		merged := make(map[string]interface{}, len(rule.Sources))
		for _, src := range rule.Sources {
			if val, ok := record[src]; ok {
				merged[src] = val
			}
		}

		record[rule.Target] = merged

		if rule.Remove {
			for _, src := range rule.Sources {
				delete(record, src)
			}
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("merge: marshal error: %w", err)
	}
	return out, nil
}
