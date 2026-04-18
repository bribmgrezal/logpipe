package dedup

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

// Deduplicator drops repeated log lines within a TTL window.
type Deduplicator struct {
	mu    sync.Mutex
	seen  map[uint64]time.Time
	ttl   time.Duration
	field string
}

// New creates a Deduplicator. ttl <= 0 defaults to 5s.
func New(ttl time.Duration, field string) *Deduplicator {
	if ttl <= 0 {
		ttl = 5 * time.Second
	}
	return &Deduplicator{
		seen:  make(map[uint64]time.Time),
		ttl:   ttl,
		field: field,
	}
}

// Allow returns true if the line is NOT a duplicate within the TTL window.
func (d *Deduplicator) Allow(line []byte) bool {
	key := d.hashLine(line)
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()
	if exp, ok := d.seen[key]; ok && now.Before(exp) {
		return false
	}
	d.seen[key] = now.Add(d.ttl)
	return true
}

// Wrap returns a middleware-style processor that skips duplicate lines.
func (d *Deduplicator) Wrap(next func([]byte) error) func([]byte) error {
	return func(line []byte) error {
		if !d.Allow(line) {
			return nil
		}
		return next(line)
	}
}

// Purge removes expired entries to free memory.
func (d *Deduplicator) Purge() {
	now := time.Now()
	d.mu.Lock()
	defer d.mu.Unlock()
	for k, exp := range d.seen {
		if now.After(exp) {
			delete(d.seen, k)
		}
	}
}

func (d *Deduplicator) hashLine(line []byte) uint64 {
	data := line
	if d.field != "" {
		var m map[string]interface{}
		if err := json.Unmarshal(line, &m); err == nil {
			if v, ok := m[d.field]; ok {
				data = []byte(fmt.Sprintf("%v", v))
			}
		}
	}
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}
