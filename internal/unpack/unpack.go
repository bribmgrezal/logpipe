package unpack

import (
	"encoding/json"
	"fmt"
)

// Rule defines a field whose string value should be unpacked as JSON
// and merged into the top-level document.
type Rule struct {
	Field  string `json:"field"`
	Remove bool   `json:"remove"`
}

// Unpacker expands a nested JSON-encoded string field into the parent object.
type Unpacker struct {
	rules []Rule
}

// New creates an Unpacker with the given rules.
func New(rules []Rule) *Unpacker {
	return &Unpacker{rules: rules}
}

// Apply parses the JSON line, unpacks configured fields, and returns
// the modified JSON. If no rules are configured the line is returned as-is.
func (u *Unpacker) Apply(line []byte) ([]byte, error) {
	if len(u.rules) == 0 {
		return line, nil
	}

	var doc map[string]interface{}
	if err := json.Unmarshal(line, &doc); err != nil {
		return nil, fmt.Errorf("unpack: invalid JSON: %w", err)
	}

	for _, rule := range u.rules {
		raw, ok := doc[rule.Field]
		if !ok {
			continue
		}
		str, ok := raw.(string)
		if !ok {
			continue
		}
		var nested map[string]interface{}
		if err := json.Unmarshal([]byte(str), &nested); err != nil {
			// field is a string but not valid JSON — skip silently
			continue
		}
		for k, v := range nested {
			doc[k] = v
		}
		if rule.Remove {
			delete(doc, rule.Field)
		}
	}

	out, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("unpack: marshal: %w", err)
	}
	return out, nil
}
