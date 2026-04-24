package rename

import (
	"encoding/json"
	"fmt"
)

// Rule defines a single field rename operation.
type Rule struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Renamer renames fields in JSON log lines.
type Renamer struct {
	rules []Rule
}

// New creates a new Renamer with the given rules.
func New(rules []Rule) *Renamer {
	return &Renamer{rules: rules}
}

// Apply renames fields in a JSON log line according to configured rules.
// Returns the modified line, or an error if the line is not valid JSON.
func (r *Renamer) Apply(line []byte) ([]byte, error) {
	if len(r.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("rename: invalid JSON: %w", err)
	}

	for _, rule := range r.rules {
		if rule.From == "" || rule.To == "" {
			continue
		}
		val, ok := record[rule.From]
		if !ok {
			continue
		}
		record[rule.To] = val
		delete(record, rule.From)
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("rename: marshal error: %w", err)
	}
	return out, nil
}
