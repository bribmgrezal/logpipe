package aggregate

import (
	"encoding/json"
	"sync"
	"time"
)

// Aggregator groups log lines by a key field and emits counts periodically.
type Aggregator struct {
	mu       sync.Mutex
	field    string
	counts   map[string]int
	interval time.Duration
	stop     chan struct{}
	out      func([]byte) error
}

// New creates an Aggregator that counts by field and flushes every interval.
func New(field string, interval time.Duration, out func([]byte) error) *Aggregator {
	if interval <= 0 {
		interval = 10 * time.Second
	}
	return &Aggregator{
		field:    field,
		counts:   make(map[string]int),
		interval: interval,
		stop:     make(chan struct{}),
		out:      out,
	}
}

// Record parses a JSON log line and increments the count for the key value.
func (a *Aggregator) Record(line []byte) error {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return err
	}
	val, ok := m[a.field]
	if !ok {
		return nil
	}
	key := toString(val)
	a.mu.Lock()
	a.counts[key]++
	a.mu.Unlock()
	return nil
}

// Start begins the flush ticker in a goroutine.
func (a *Aggregator) Start() {
	go func() {
		ticker := time.NewTicker(a.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.Flush()
			case <-a.stop:
				return
			}
		}
	}()
}

// Stop halts the background flush goroutine.
func (a *Aggregator) Stop() { close(a.stop) }

// Flush emits the current counts as a JSON line and resets.
func (a *Aggregator) Flush() {
	a.mu.Lock()
	snap := make(map[string]int, len(a.counts))
	for k, v := range a.counts {
		snap[k] = v
	}
	a.counts = make(map[string]int)
	a.mu.Unlock()
	if len(snap) == 0 {
		return
	}
	out := map[string]interface{}{"_aggregate": a.field, "counts": snap}
	b, _ := json.Marshal(out)
	_ = a.out(b)
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	b, _ := json.Marshal(v)
	return string(b)
}
