package validate

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Rule defines a single field validation rule.
type Rule struct {
	Field   string `json:"field"`
	Pattern string `json:"pattern"`
	MinLen  int    `json:"min_len"`
	MaxLen  int    `json:"max_len"`
}

// Validator applies regex and length constraints to JSON log fields.
type Validator struct {
	rules   []Rule
	compiled []*regexp.Regexp
}

// New creates a Validator from the given rules.
func New(rules []Rule) (*Validator, error) {
	v := &Validator{rules: rules}
	for _, r := range rules {
		if r.Pattern == "" {
			v.compiled = append(v.compiled, nil)
			continue
		}
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, fmt.Errorf("validate: invalid pattern %q: %w", r.Pattern, err)
		}
		v.compiled = append(v.compiled, re)
	}
	return v, nil
}

// Check validates a JSON line against all rules.
// Returns an error describing the first violation, or nil if all pass.
func (v *Validator) Check(line []byte) error {
	if len(v.rules) == 0 {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return fmt.Errorf("validate: invalid JSON: %w", err)
	}
	for i, r := range v.rules {
		val, ok := m[r.Field]
		if !ok {
			return fmt.Errorf("validate: missing field %q", r.Field)
		}
		s := fmt.Sprintf("%v", val)
		if r.MinLen > 0 && len(s) < r.MinLen {
			return fmt.Errorf("validate: field %q length %d below min %d", r.Field, len(s), r.MinLen)
		}
		if r.MaxLen > 0 && len(s) > r.MaxLen {
			return fmt.Errorf("validate: field %q length %d exceeds max %d", r.Field, len(s), r.MaxLen)
		}
		if re := v.compiled[i]; re != nil && !re.MatchString(s) {
			return fmt.Errorf("validate: field %q value %q does not match pattern %q", r.Field, s, r.Pattern)
		}
	}
	return nil
}
