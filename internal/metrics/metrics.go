package metrics

import (
	"fmt"
	"sync/atomic"
)

// Counter tracks pipeline processing statistics.
type Counter struct {
	Received  atomic.Int64
	Passed    atomic.Int64
	Filtered  atomic.Int64
	Errors    atomic.Int64
}

// Global is the default shared counter.
var Global = &Counter{}

// IncReceived increments the received lines counter.
func (c *Counter) IncReceived() { c.Received.Add(1) }

// IncPassed increments the passed (not filtered) lines counter.
func (c *Counter) IncPassed() { c.Passed.Add(1) }

// IncFiltered increments the filtered lines counter.
func (c *Counter) IncFiltered() { c.Filtered.Add(1) }

// IncErrors increments the error counter.
func (c *Counter) IncErrors() { c.Errors.Add(1) }

// Snapshot returns a point-in-time copy of all counters.
type Snapshot struct {
	Received int64
	Passed   int64
	Filtered int64
	Errors   int64
}

// String returns a human-readable summary of the snapshot.
func (s Snapshot) String() string {
	return fmt.Sprintf("received=%d passed=%d filtered=%d errors=%d",
		s.Received, s.Passed, s.Filtered, s.Errors)
}

// Snapshot captures current counter values.
func (c *Counter) Snapshot() Snapshot {
	return Snapshot{
		Received: c.Received.Load(),
		Passed:   c.Passed.Load(),
		Filtered: c.Filtered.Load(),
		Errors:   c.Errors.Load(),
	}
}

// Reset zeroes all counters.
func (c *Counter) Reset() {
	c.Received.Store(0)
	c.Passed.Store(0)
	c.Filtered.Store(0)
	c.Errors.Store(0)
}
