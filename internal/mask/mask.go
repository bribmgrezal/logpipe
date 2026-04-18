package mask

import (
	"encoding/json"
	"regexp"
	"strings"
)

// Rule defines a masking rule for a field using a regex pattern.
type Rule struct {
	Field   string `json:"field"`
	Pattern string `json:"pattern"`
	Mask    string `json:"mask"`
}

// Masker applies regex-based masking to JSON log lines.
type Masker struct {
	rules []compiledRule
}

type compiledRule struct {
	field string
	re    *regexp.Regexp
	mask  string
}

// New creates a Masker from the given rules.
func New(rules []Rule) (*Masker, error) {
	var compiled []compiledRule
	for _, r := range rules {
		re, err := regexp.Compile(r.Pattern)
		if err != nil {
			return nil, err
		}
		mask := r.Mask
		if mask == "" {
			mask = "***"
		}
		compiled = append(compiled, compiledRule{field: r.Field, re: re, mask: mask})
	}
	return &Masker{rules: compiled}, nil
}

// Apply masks matching patterns in the specified fields of a JSON line.
func (m *Masker) Apply(line string) (string, error) {
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return line, err
	}
	for _, r := range m.rules {
		if val, ok := obj[r.field]; ok {
			if s, ok := val.(string); ok {
				obj[r.field] = r.re.ReplaceAllString(s, r.mask)
			}
		}
	}
	out, err := json.Marshal(obj)
	if err != nil {
		return line, err
	}
	return strings.TrimSpace(string(out)), nil
}
