package enrich

import (
	"encoding/json"
	"fmt"
)

// Rule defines a static field to add or a field to copy/rename.
type Rule struct {
	Field  string `json:"field"`
	Value  string `json:"value,omitempty"`
	CopyOf string `json:"copy_of,omitempty"`
}

// Enricher applies enrichment rules to log lines.
type Enricher struct {
	rules []Rule
}

// New creates an Enricher with the given rules.
func New(rules []Rule) *Enricher {
	return &Enricher{rules: rules}
}

// Apply enriches a JSON log line and returns the modified JSON.
func (e *Enricher) Apply(line []byte) ([]byte, error) {
	if len(e.rules) == 0 {
		return line, nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("enrich: invalid json: %w", err)
	}
	for _, r := range e.rules {
		if r.CopyOf != "" {
			if v, ok := m[r.CopyOf]; ok {
				m[r.Field] = v
			}
		} else {
			m[r.Field] = r.Value
		}
	}
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("enrich: marshal: %w", err)
	}
	return out, nil
}

// Wrap returns a middleware function that enriches each line.
func (e *Enricher) Wrap(next func([]byte) error) func([]byte) error {
	return func(line []byte) error {
		enriched, err := e.Apply(line)
		if err != nil {
			return err
		}
		return next(enriched)
	}
}
