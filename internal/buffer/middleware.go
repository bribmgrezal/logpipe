package buffer

import (
	"encoding/json"
)

// Middleware wraps a RingBuffer and provides buffered write functionality
// for use in the pipeline, storing raw JSON log lines.
type Middleware struct {
	buf *RingBuffer
}

// NewMiddleware creates a Middleware backed by a RingBuffer of given capacity.
func NewMiddleware(capacity int) *Middleware {
	return &Middleware{buf: New(capacity)}
}

// Write stores the raw line in the ring buffer and returns it unchanged.
// Returns an error if the line is not valid JSON.
func (m *Middleware) Write(line string) (string, error) {
	if !json.Valid([]byte(line)) {
		return "", &InvalidJSONError{Line: line}
	}
	m.buf.Push(line)
	return line, nil
}

// Snapshot returns all buffered lines.
func (m *Middleware) Snapshot() []string {
	return m.buf.Snapshot()
}

// Len returns the number of buffered lines.
func (m *Middleware) Len() int {
	return m.buf.Len()
}

// Reset clears the buffer.
func (m *Middleware) Reset() {
	m.buf.Reset()
}

// InvalidJSONError is returned when a non-JSON line is written.
type InvalidJSONError struct {
	Line string
}

func (e *InvalidJSONError) Error() string {
	return "buffer: invalid JSON: " + e.Line
}
