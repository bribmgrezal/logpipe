package dedup

import (
	"encoding/json"
	"testing"
	"time"
)

func encode(m map[string]interface{}) []byte {
	b, _ := json.Marshal(m)
	return b
}

func TestAllow_FirstLineAllowed(t *testing.T) {
	d := New(time.Second, "")
	line := []byte(`{"msg":"hello"}`)
	if !d.Allow(line) {
		t.Fatal("expected first line to be allowed")
	}
}

func TestAllow_DuplicateBlocked(t *testing.T) {
	d := New(time.Second, "")
	line := []byte(`{"msg":"hello"}`)
	d.Allow(line)
	if d.Allow(line) {
		t.Fatal("expected duplicate to be blocked")
	}
}

func TestAllow_DifferentLinesAllowed(t *testing.T) {
	d := New(time.Second, "")
	a := []byte(`{"msg":"hello"}`)
	b := []byte(`{"msg":"world"}`)
	d.Allow(a)
	if !d.Allow(b) {
		t.Fatal("expected different line to be allowed")
	}
}

func TestAllow_FieldDedup(t *testing.T) {
	d := New(time.Second, "msg")
	a := encode(map[string]interface{}{"msg": "hello", "ts": 1})
	b := encode(map[string]interface{}{"msg": "hello", "ts": 2})
	d.Allow(a)
	if d.Allow(b) {
		t.Fatal("expected field-based duplicate to be blocked")
	}
}

func TestPurge_RemovesExpired(t *testing.T) {
	d := New(10*time.Millisecond, "")
	line := []byte(`{"msg":"purge"}`)
	d.Allow(line)
	time.Sleep(20 * time.Millisecond)
	d.Purge()
	d.mu.Lock()
	l := len(d.seen)
	d.mu.Unlock()
	if l != 0 {
		t.Fatalf("expected 0 entries after purge, got %d", l)
	}
}

func TestWrap_SkipsDuplicates(t *testing.T) {
	d := New(time.Second, "")
	count := 0
	next := func(line []byte) error { count++; return nil }
	processor := d.Wrap(next)
	line := []byte(`{"msg":"dup"}`)
	processor(line)
	processor(line)
	processor(line)
	if count != 1 {
		t.Fatalf("expected 1 call, got %d", count)
	}
}

func TestNew_DefaultTTL(t *testing.T) {
	d := New(0, "")
	if d.ttl != 5*time.Second {
		t.Fatalf("expected default ttl 5s, got %v", d.ttl)
	}
}
