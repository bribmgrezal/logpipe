package ratelimit

import (
	"sync"
	"time"
)

// Limiter is a token-bucket rate limiter for log lines.
type Limiter struct {
	mu       sync.Mutex
	rate     int // tokens per second
	tokens   int
	max      int
	lastTick time.Time
}

// New creates a Limiter that allows up to rate events per second.
func New(rate int) *Limiter {
	if rate <= 0 {
		rate = 1
	}
	return &Limiter{
		rate:     rate,
		tokens:   rate,
		max:      rate,
		lastTick: time.Now(),
	}
}

// Allow returns true if the event is permitted under the rate limit.
func (l *Limiter) Allow() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.refill()
	if l.tokens > 0 {
		l.tokens--
		return true
	}
	return false
}

// refill adds tokens based on elapsed time since last call.
func (l *Limiter) refill() {
	now := time.Now()
	elapsed := now.Sub(l.lastTick).Seconds()
	add := int(elapsed * float64(l.rate))
	if add > 0 {
		l.tokens += add
		if l.tokens > l.max {
			l.tokens = l.max
		}
		l.lastTick = now
	}
}

// Wrap returns a filtered channel that drops lines exceeding the rate limit.
func Wrap(in <-chan string, rate int) <-chan string {
	out := make(chan string, 64)
	limiter := New(rate)
	go func() {
		defer close(out)
		for line := range in {
			if limiter.Allow() {
				out <- line
			}
		}
	}()
	return out
}
