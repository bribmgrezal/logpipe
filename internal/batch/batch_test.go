package batch

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func encode(v map[string]any) []byte {
	b, _ := json.Marshal(v)
	return b
}

func TestWrite_BuffersLines(t *testing.T) {
	var mu sync.Mutex
	var got [][]map[string]any
	b := New(10, time.Minute, func(batch []map[string]any) {
		mu.Lock()
		got = append(got, batch)
		mu.Unlock()
	})
	defer b.Stop()
	if err := b.Write(encode(map[string]any{"msg": "hello"})); err != nil {
		t.Fatal(err)
	}
	b.Flush()
	mu.Lock()
	defer mu.Unlock()
	if len(got) != 1 || len(got[0]) != 1 {
		t.Fatalf("expected 1 batch with 1 entry, got %v", got)
	}
}

func TestWrite_FlushesOnSize(t *testing.T) {
	var mu sync.Mutex
	var count int
	b := New(2, time.Minute, func(batch []map[string]any) {
		mu.Lock()
		count += len(batch)
		mu.Unlock()
	})
	defer b.Stop()
	b.Write(encode(map[string]any{"a": "1"}))
	b.Write(encode(map[string]any{"a": "2"}))
	time.Sleep(10 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if count != 2 {
		t.Fatalf("expected 2 flushed, got %d", count)
	}
}

func TestWrite_InvalidJSON(t *testing.T) {
	b := New(10, time.Minute, func(_ []map[string]any) {})
	defer b.Stop()
	if err := b.Write([]byte("not-json")); err == nil {
		t.Fatal("expected error")
	}
}

func TestFlush_EmptyNoop(t *testing.T) {
	called := false
	b := New(10, time.Minute, func(_ []map[string]any) { called = true })
	defer b.Stop()
	b.Flush()
	if called {
		t.Fatal("flush on empty should not call flushFn")
	}
}

func TestTicker_AutoFlush(t *testing.T) {
	var mu sync.Mutex
	var flushed bool
	b := New(100, 30*time.Millisecond, func(_ []map[string]any) {
		mu.Lock()
		flushed = true
		mu.Unlock()
	})
	defer b.Stop()
	b.Write(encode(map[string]any{"x": "y"}))
	time.Sleep(80 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	if !flushed {
		t.Fatal("expected auto-flush via ticker")
	}
}
