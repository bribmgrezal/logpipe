package schema

import (
	"encoding/json"
	"fmt"
)

// Rule defines a field validation rule.
type Rule struct {
	Field    string `json:"field"`
	Required bool   `json:"required"`
	Type     string `json:"type"` // string, number, bool
}

// Validator validates JSON log lines against a set of rules.
type Validator struct {
	rules []Rule
}

// New creates a new Validator.
func New(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// Validate checks a JSON line against all rules.
// Returns an error if any rule is violated, nil otherwise.
func (v *Validator) Validate(line []byte) error {
	if len(v.rules) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return fmt.Errorf("invalid json: %w", err)
	}
	for _, r := range v.rules {
		val, exists := m[r.Field]
		if r.Required && !exists {
			return fmt.Errorf("missing required field: %s", r.Field)
		}
		if !exists {
			continue
		}
		if r.Type != "" {
			if err := checkType(r.Field, r.Type, val); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkType(field, expected string, val interface{}) error {
	switch expected {
	case "string":
		if _, ok := val.(string); !ok {
			return fmt.Errorf("field %s: expected string", field)
		}
	case "number":
		if _, ok := val.(float64); !ok {
			return fmt.Errorf("field %s: expected number", field)
		}
	case "bool":
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("field %s: expected bool", field)
		}
	}
	return nil
}
