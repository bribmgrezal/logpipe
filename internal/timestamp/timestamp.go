package timestamp

import (
	"encoding/json"
	"fmt"
	"time"
)

// Processor rewrites a timestamp field into a normalized format.
type Processor struct {
	field  string
	inFmt  string
	outFmt string
}

// New creates a Processor. inFmt/outFmt use Go time layout strings.
// If outFmt is empty, time.RFC3339 is used.
func New(field, inFmt, outFmt string) (*Processor, error) {
	if field == "" {
		return nil, fmt.Errorf("timestamp: field must not be empty")
	}
	if outFmt == "" {
		outFmt = time.RFC3339
	}
	return &Processor{field: field, inFmt: inFmt, outFmt: outFmt}, nil
}

// Apply parses the timestamp field in line and rewrites it.
// Returns the original line unchanged on any error.
func (p *Processor) Apply(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("timestamp: invalid JSON: %w", err)
	}

	raw, ok := m[p.field]
	if !ok {
		return line, nil
	}

	str, ok := raw.(string)
	if !ok {
		return line, nil
	}

	var t time.Time
	var err error
	if p.inFmt == "" {
		t, err = time.Parse(time.RFC3339, str)
	} else {
		t, err = time.Parse(p.inFmt, str)
	}
	if err != nil {
		return line, nil
	}

	m[p.field] = t.UTC().Format(p.outFmt)
	out, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("timestamp: marshal error: %w", err)
	}
	return out, nil
}
