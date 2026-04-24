package extract

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Rule defines a single extraction rule: parse a source field's string value
// using a delimiter and write the resulting parts into named target fields.
type Rule struct {
	Field     string   `json:"field"`
	Delimiter string   `json:"delimiter"`
	Targets   []string `json:"targets"`
}

// Extractor applies extraction rules to JSON log lines.
type Extractor struct {
	rules []Rule
}

// New creates an Extractor with the given rules.
func New(rules []Rule) *Extractor {
	return &Extractor{rules: rules}
}

// Apply parses line as JSON, applies all extraction rules, and returns the
// modified JSON. Lines that are not valid JSON are returned unchanged with an
// error.
func (e *Extractor) Apply(line string) (string, error) {
	if len(e.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return line, fmt.Errorf("extract: invalid JSON: %w", err)
	}

	for _, r := range e.rules {
		delim := r.Delimiter
		if delim == "" {
			delim = " "
		}

		raw, ok := record[r.Field]
		if !ok {
			continue
		}
		str, ok := raw.(string)
		if !ok {
			continue
		}

		parts := strings.SplitN(str, delim, len(r.Targets))
		for i, target := range r.Targets {
			if i < len(parts) {
				record[target] = parts[i]
			}
		}
	}

	out, err := json.Marshal(record)
	if err != nil {
		return line, fmt.Errorf("extract: marshal error: %w", err)
	}
	return string(out), nil
}
