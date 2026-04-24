package dedupe

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
)

// Deduper removes duplicate log lines based on a configurable window and optional field key.
type Deduper struct {
	mu      sync.Mutex
	seen    map[string]struct{}
	fields  []string
	maxSize int
}

// New creates a new Deduper. If fields is non-empty, deduplication is based on
// the values of those fields only; otherwise the entire line is hashed.
func New(fields []string, maxSize int) *Deduper {
	if maxSize <= 0 {
		maxSize = 10000
	}
	return &Deduper{
		seen:    make(map[string]struct{}, maxSize),
		fields:  fields,
		maxSize: maxSize,
	}
}

// Allow returns true if the line has not been seen before.
func (d *Deduper) Allow(line []byte) (bool, error) {
	key, err := d.computeKey(line)
	if err != nil {
		return false, err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, exists := d.seen[key]; exists {
		return false, nil
	}
	if len(d.seen) >= d.maxSize {
		d.seen = make(map[string]struct{}, d.maxSize)
	}
	d.seen[key] = struct{}{}
	return true, nil
}

// Reset clears the seen set.
func (d *Deduper) Reset() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen = make(map[string]struct{}, d.maxSize)
}

func (d *Deduper) computeKey(line []byte) (string, error) {
	if len(d.fields) == 0 {
		h := sha256.Sum256(line)
		return fmt.Sprintf("%x", h), nil
	}
	var record map[string]interface{}
	if err := json.Unmarshal(line, &record); err != nil {
		return "", fmt.Errorf("dedupe: invalid JSON: %w", err)
	}
	parts := make(map[string]interface{}, len(d.fields))
	for _, f := range d.fields {
		parts[f] = record[f]
	}
	b, err := json.Marshal(parts)
	if err != nil {
		return "", err
	}
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h), nil
}
