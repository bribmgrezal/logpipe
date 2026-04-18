package throttle

import (
	"encoding/json"
	"sync"
	"time"
)

// Throttler suppresses repeated identical log lines within a window.
type Throttler struct {
	mu      sync.Mutex
	window  time.Duration
	seen    map[string]time.Time
	clock   func() time.Time
}

// New creates a Throttler with the given suppression window.
func New(window time.Duration) *Throttler {
	if window <= 0 {
		window = time.Second
	}
	return &Throttler{
		window: window,
		seen:   make(map[string]time.Time),
		clock:  time.Now,
	}
}

// Allow returns true if the line should be forwarded (not throttled).
func (t *Throttler) Allow(line []byte) bool {
	key := string(line)
	t.mu.Lock()
	defer t.mu.Unlock()
	now := t.clock()
	if last, ok := t.seen[key]; ok && now.Sub(last) < t.window {
		return false
	}
	t.seen[key] = now
	return true
}

// Reset clears the suppression state.
func (t *Throttler) Reset() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.seen = make(map[string]time.Time)
}

// Wrap returns a middleware function that throttles lines.
func Wrap(t *Throttler, next func([]byte) error) func([]byte) error {
	return func(line []byte) error {
		if !json.Valid(line) {
			return next(line)
		}
		if t.Allow(line) {
			return next(line)
		}
		return nil
	}
}
