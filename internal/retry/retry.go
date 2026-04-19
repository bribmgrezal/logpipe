package retry

import (
	"encoding/json"
	"time"
)

// Retryer holds retry configuration and state.
type Retryer struct {
	maxAttempts int
	delay       time.Duration
	writer      func([]byte) error
}

// New creates a Retryer with the given max attempts and delay between retries.
func New(maxAttempts int, delay time.Duration, writer func([]byte) error) *Retryer {
	if maxAttempts < 1 {
		maxAttempts = 1
	}
	return &Retryer{
		maxAttempts: maxAttempts,
		delay:       delay,
		writer:      writer,
	}
}

// Write attempts to write the JSON line, retrying on failure up to maxAttempts times.
func (r *Retryer) Write(line []byte) error {
	if !json.Valid(line) {
		return &InvalidJSONError{line: string(line)}
	}
	var err error
	for i := 0; i < r.maxAttempts; i++ {
		if err = r.writer(line); err == nil {
			return nil
		}
		if i < r.maxAttempts-1 && r.delay > 0 {
			time.Sleep(r.delay)
		}
	}
	return err
}

// InvalidJSONError is returned when the input is not valid JSON.
type InvalidJSONError struct {
	line string
}

func (e *InvalidJSONError) Error() string {
	return "retry: invalid JSON: " + e.line
}
