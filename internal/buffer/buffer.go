package buffer

import (
	"sync"
)

// RingBuffer holds a fixed number of recent log lines in memory.
type RingBuffer struct {
	mu      sync.Mutex
	data    []string
	cap     int
	head    int
	count   int
}

// New creates a RingBuffer with the given capacity.
func New(capacity int) *RingBuffer {
	if capacity <= 0 {
		capacity = 100
	}
	return &RingBuffer{
		data: make([]string, capacity),
		cap:  capacity,
	}
}

// Push adds a line to the buffer, overwriting the oldest entry when full.
func (r *RingBuffer) Push(line string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.data[r.head] = line
	r.head = (r.head + 1) % r.cap
	if r.count < r.cap {
		r.count++
	}
}

// Snapshot returns a copy of all buffered lines in insertion order.
func (r *RingBuffer) Snapshot() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.count == 0 {
		return nil
	}
	out := make([]string, r.count)
	start := (r.head - r.count + r.cap) % r.cap
	for i := 0; i < r.count; i++ {
		out[i] = r.data[(start+i)%r.cap]
	}
	return out
}

// Len returns the current number of items in the buffer.
func (r *RingBuffer) Len() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.count
}

// Reset clears the buffer.
func (r *RingBuffer) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.head = 0
	r.count = 0
}
