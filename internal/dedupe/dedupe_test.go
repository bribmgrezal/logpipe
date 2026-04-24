package dedupe

import (
	"encoding/json"
	"testing"
)

func encode(t *testing.T, v map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	return b
}

func TestAllow_UniqueLines(t *testing.T) {
	d := New(nil, 100)
	line := []byte(`{"level":"info","msg":"hello"}`)
	ok, err := d.Allow(line)
	if err != nil || !ok {
		t.Fatalf("expected first line allowed, got ok=%v err=%v", ok, err)
	}
}

func TestAllow_DuplicateBlocked(t *testing.T) {
	d := New(nil, 100)
	line := []byte(`{"level":"info","msg":"hello"}`)
	d.Allow(line) //nolint
	ok, err := d.Allow(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected duplicate to be blocked")
	}
}

func TestAllow_DifferentLinesAllowed(t *testing.T) {
	d := New(nil, 100)
	a := []byte(`{"msg":"one"}`)
	b := []byte(`{"msg":"two"}`)
	d.Allow(a) //nolint
	ok, err := d.Allow(b)
	if err != nil || !ok {
		t.Fatalf("expected different line allowed, got ok=%v err=%v", ok, err)
	}
}

func TestAllow_FieldDedup(t *testing.T) {
	d := New([]string{"id"}, 100)
	a := encode(t, map[string]interface{}{"id": "abc", "ts": "1"})
	b := encode(t, map[string]interface{}{"id": "abc", "ts": "2"})
	d.Allow(a) //nolint
	ok, err := d.Allow(b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected field-based duplicate to be blocked")
	}
}

func TestAllow_InvalidJSON(t *testing.T) {
	d := New([]string{"id"}, 100)
	_, err := d.Allow([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestReset_ClearsState(t *testing.T) {
	d := New(nil, 100)
	line := []byte(`{"msg":"hi"}`)
	d.Allow(line) //nolint
	d.Reset()
	ok, err := d.Allow(line)
	if err != nil || !ok {
		t.Fatalf("expected line allowed after reset, got ok=%v err=%v", ok, err)
	}
}

func TestNew_DefaultMaxSize(t *testing.T) {
	d := New(nil, 0)
	if d.maxSize != 10000 {
		t.Fatalf("expected default maxSize 10000, got %d", d.maxSize)
	}
}
