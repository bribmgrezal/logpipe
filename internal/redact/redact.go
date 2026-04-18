package redact

import (
	"encoding/json"
	"strings"
)

// Rule defines a field to redact and how to redact it.
type Rule struct {
	Field   string `json:"field"`
	Replace string `json:"replace"` // default: "***"
}

// Redactor applies redaction rules to JSON log lines.
type Redactor struct {
	rules []Rule
}

// New creates a Redactor with the given rules.
func New(rules []Rule) *Redactor {
	for i := range rules {
		if rules[i].Replace == "" {
			rules[i].Replace = "***"
		}
	}
	return &Redactor{rules: rules}
}

// Apply redacts sensitive fields from a JSON line.
// Returns the modified line or the original on error.
func (r *Redactor) Apply(line string) string {
	if len(r.rules) == 0 {
		return line
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line
	}

	for _, rule := range r.rules {
		redactNested(obj, rule.Field, rule.Replace)
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return line
	}
	return string(out)
}

// redactNested supports dot-notation for nested fields.
func redactNested(obj map[string]interface{}, field, replace string) {
	parts := strings.SplitN(field, ".", 2)
	if len(parts) == 1 {
		if _, ok := obj[field]; ok {
			obj[field] = replace
		}
		return
	}
	if nested, ok := obj[parts[0]].(map[string]interface{}); ok {
		redactNested(nested, parts[1], replace)
	}
}
