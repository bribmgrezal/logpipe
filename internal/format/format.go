package format

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Formatter converts a raw JSON log line into a formatted string.
type Formatter struct {
	template string
	fields   []string
}

// New creates a Formatter. If template is empty, pretty JSON is used.
func New(template string) *Formatter {
	fields := extractFields(template)
	return &Formatter{template: template, fields: fields}
}

// Apply formats a JSON log line according to the template.
func (f *Formatter) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("format: invalid JSON: %w", err)
	}
	if f.template == "" {
		out, err := json.MarshalIndent(m, "", "  ")
		if err != nil {
			return nil, err
		}
		return out, nil
	}
	result := f.template
	for _, field := range f.fields {
		val := ""
		if v, ok := m[field]; ok {
			val = fmt.Sprintf("%v", v)
		}
		result = strings.ReplaceAll(result, "{"+field+"}", val)
	}
	return []byte(result), nil
}

func extractFields(tmpl string) []string {
	var fields []string
	for {
		start := strings.Index(tmpl, "{")
		if start == -1 {
			break
		}
		end := strings.Index(tmpl[start:], "}")
		if end == -1 {
			break
		}
		fields = append(fields, tmpl[start+1:start+end])
		tmpl = tmpl[start+end+1:]
	}
	return fields
}
