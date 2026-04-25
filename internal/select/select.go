package select_

import (
	"encoding/json"
	"fmt"
)

// Selector retains only specified fields from each log line.
type Selector struct {
	fields []string
}

// New creates a Selector that keeps only the given fields.
// If fields is empty, all fields are retained.
func New(fields []string) *Selector {
	return &Selector{fields: fields}
}

// Apply filters a JSON log line, retaining only the configured fields.
// If no fields are configured, the line is returned unchanged.
// Returns an error if the line is not valid JSON.
func (s *Selector) Apply(line []byte) ([]byte, error) {
	if len(s.fields) == 0 {
		return line, nil
	}

	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return nil, fmt.Errorf("select: invalid JSON: %w", err)
	}

	out := make(map[string]interface{}, len(s.fields))
	for _, f := range s.fields {
		if v, ok := record[f]; ok {
			out[f] = v
		}
	}

	result, err := json.Marshal(out)
	if err != nil {
		return nil, fmt.Errorf("select: marshal error: %w", err)
	}
	return result, nil
}
