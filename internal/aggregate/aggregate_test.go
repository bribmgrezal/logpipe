package aggregate

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestRecord_ValidLine(t *testing.T) {
	a := New("level", time.Second, func([]byte) error { return nil })
	if err := a.Record([]byte(`{"level":"info","msg":"ok"}`)); err != nil {
		t.Fatal(err)
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.counts["info"] != 1 {
		t.Fatalf("expected 1, got %d", a.counts["info"])
	}
}

func TestRecord_InvalidJSON(t *testing.T) {
	a := New("level", time.Second, func([]byte) error { return nil })
	if err := a.Record([]byte(`not-json`)); err == nil {
		t.Fatal("expected error")
	}
}

func TestRecord_MissingField(t *testing.T) {
	a := New("level", time.Second, func([]byte) error { return nil })
	_ = a.Record([]byte(`{"msg":"no level"}`)) 
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.counts) != 0 {
		t.Fatal("expected empty counts")
	}
}

func TestFlush_EmitsAndResets(t *testing.T) {
	var mu sync.Mutex
	var got []map[string]interface{}
	a := New("level", time.Second, func(b []byte) error {
		var m map[string]interface{}
		_ = json.Unmarshal(b, &m)
		mu.Lock()
		got = append(got, m)
		mu.Unlock()
		return nil
	})
	_ = a.Record([]byte(`{"level":"warn"}`))
	_ = a.Record([]byte(`{"level":"warn"}`))
	_ = a.Record([]byte(`{"level":"info"}`))
	a.Flush()
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 {
		t.Fatalf("expected 1 flush, got %d", len(got))
	}
	if got[0]["_aggregate"] != "level" {
		t.Fatal("missing _aggregate key")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.counts) != 0 {
		t.Fatal("counts not reset after flush")
	}
}

func TestFlush_EmptyNoOutput(t *testing.T) {
	called := false
	a := New("level", time.Second, func([]byte) error { called = true; return nil })
	a.Flush()
	if called {
		t.Fatal("should not emit on empty counts")
	}
}
