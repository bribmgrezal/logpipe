package batch

import (
	"encoding/json"
	"sync"
	"time"
)

// Batcher accumulates log lines and flushes them as a JSON array.
type Batcher struct {
	mu       sync.Mutex
	buf      []map[string]any
	size     int
	interval time.Duration
	flushFn  func([]map[string]any)
	stop     chan struct{}
}

// New creates a Batcher that flushes when size is reached or interval elapses.
func New(size int, interval time.Duration, flushFn func([]map[string]any)) *Batcher {
	if size <= 0 {
		size = 100
	}
	if interval <= 0 {
		interval = 5 * time.Second
	}
	b := &Batcher{
		size:     size,
		interval: interval,
		flushFn:  flushFn,
		stop:     make(chan struct{}),
	}
	go b.ticker()
	return b
}

// Write accepts a JSON log line and buffers it.
func (b *Batcher) Write(line []byte) error {
	var entry map[string]any
	if err := json.Unmarshal(line, &entry); err != nil {
		return err
	}
	b.mu.Lock()
	b.buf = append(b.buf, entry)
	should := len(b.buf) >= b.size
	b.mu.Unlock()
	if should {
		b.Flush()
	}
	return nil
}

// Flush emits the current buffer and resets it.
func (b *Batcher) Flush() {
	b.mu.Lock()
	if len(b.buf) == 0 {
		b.mu.Unlock()
		return
	}
	out := make([]map[string]any, len(b.buf))
	copy(out, b.buf)
	b.buf = b.buf[:0]
	b.mu.Unlock()
	b.flushFn(out)
}

// Stop halts the background ticker.
func (b *Batcher) Stop() {
	close(b.stop)
	b.Flush()
}

func (b *Batcher) ticker() {
	t := time.NewTicker(b.interval)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			b.Flush()
		case <-b.stop:
			return
		}
	}
}
