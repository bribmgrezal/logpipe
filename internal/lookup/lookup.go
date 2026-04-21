package lookup

import (
	"encoding/json"
	"fmt"
)

// Rule defines a lookup enrichment: if field equals value, add fields from the table.
type Rule struct {
	Field  string            `json:"field"`
	Table  map[string]string `json:"table"`
	Target string            `json:"target"` // optional prefix for added keys
}

// Lookup enriches log lines by looking up field values in static tables.
type Lookup struct {
	rules []Rule
}

// New creates a new Lookup with the given rules.
func New(rules []Rule) *Lookup {
	return &Lookup{rules: rules}
}

// Apply enriches a JSON log line using lookup rules.
// For each rule, if the field exists in the line, its value is used as a key
// to look up an enrichment value which is added to the output.
func (l *Lookup) Apply(line []byte) ([]byte, error) {
	if len(l.rules) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("lookup: invalid JSON: %w", err)
	}

	for _, rule := range l.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		key := fmt.Sprintf("%v", val)
		enriched, found := rule.Table[key]
		if !found {
			continue
		}
		target := rule.Field + "_lookup"
		if rule.Target != "" {
			target = rule.Target
		}
		record[target] = enriched
	}

	out, err := json.Marshal(record)
	if err != nil {
		return nil, fmt.Errorf("lookup: marshal error: %w", err)
	}
	return out, nil
}
