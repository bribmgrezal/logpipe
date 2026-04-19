package split

import (
	"encoding/json"
	"fmt"
)

// Rule defines a field and its expected value to route to a named output.
type Rule struct {
	Field  string `json:"field"`
	Value  string `json:"value"`
	Output string `json:"output"`
}

// Splitter fans out log lines to named output buckets based on field values.
type Splitter struct {
	rules    []Rule
	fallback string
}

// New creates a Splitter with the given rules and optional fallback output name.
func New(rules []Rule, fallback string) *Splitter {
	return &Splitter{rules: rules, fallback: fallback}
}

// Apply returns the output name for the given JSON log line.
// Returns fallback if no rule matches. Returns error on invalid JSON.
func (s *Splitter) Apply(line []byte) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return "", fmt.Errorf("split: invalid JSON: %w", err)
	}

	for _, r := range s.rules {
		v, ok := obj[r.Field]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", v) == r.Value {
			return r.Output, nil
		}
	}

	return s.fallback, nil
}
