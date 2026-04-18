package ratelimit

import (
	"testing"
	"time"
)

func TestNew_DefaultsInvalidRate(t *testing.T) {
	l := New(0)
	if l.rate != 1 {
		t.Fatalf("expected rate=1, got %d", l.rate)
	}
}

func TestAllow_WithinLimit(t *testing.T) {
	l := New(5)
	for i := 0; i < 5; i++ {
		if !l.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	l := New(3)
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("expected Allow()=false after exhausting tokens")
	}
}

func TestAllow_RefillsOverTime(t *testing.T) {
	l := New(2)
	l.Allow()
	l.Allow()
	if l.Allow() {
		t.Fatal("should be empty before refill")
	}
	time.Sleep(1100 * time.Millisecond)
	if !l.Allow() {
		t.Fatal("expected token after refill")
	}
}

func TestWrap_DropsOverLimit(t *testing.T) {
	in := make(chan string, 20)
	for i := 0; i < 20; i++ {
		in <- `{"msg":"hello"}`
	}
	close(in)

	out := Wrap(in, 5)
	count := 0
	for range out {
		count++
	}
	if count > 5 {
		t.Fatalf("expected at most 5 lines, got %d", count)
	}
}
