package sampler

import (
	"encoding/json"
	"sync/atomic"
)

// Sampler drops log lines based on a 1-in-N sampling rate.
type Sampler struct {
	rate    uint64
	counter atomic.Uint64
}

// New creates a Sampler that keeps every nth message.
// A rate of 0 or 1 keeps all messages.
func New(rate uint64) *Sampler {
	if rate == 0 {
		rate = 1
	}
	return &Sampler{rate: rate}
}

// Allow returns true if the current message should be forwarded.
func (s *Sampler) Allow() bool {
	n := s.counter.Add(1)
	return n%s.rate == 1
}

// Reset resets the internal counter.
func (s *Sampler) Reset() {
	s.counter.Store(0)
}

// Wrap wraps a writer function, sampling incoming JSON lines.
func (s *Sampler) Wrap(next func([]byte) error) func([]byte) error {
	return func(line []byte) error {
		if !json.Valid(line) {
			return nil
		}
		if !s.Allow() {
			return nil
		}
		return next(line)
	}
}
