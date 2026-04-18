package buffer

import (
	"testing"
)

func TestMiddleware_WriteValid(t *testing.T) {
	m := NewMiddleware(10)
	out, err := m.Write(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != `{"level":"info","msg":"hello"}` {
		t.Fatalf("unexpected output: %s", out)
	}
	if m.Len() != 1 {
		t.Fatalf("expected 1 buffered line, got %d", m.Len())
	}
}

func TestMiddleware_WriteInvalidJSON(t *testing.T) {
	m := NewMiddleware(10)
	_, err := m.Write("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	if m.Len() != 0 {
		t.Fatal("invalid line should not be buffered")
	}
}

func TestMiddleware_Snapshot(t *testing.T) {
	m := NewMiddleware(5)
	m.Write(`{"a":1}`)
	m.Write(`{"b":2}`)
	snap := m.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 lines in snapshot, got %d", len(snap))
	}
	if snap[0] != `{"a":1}` || snap[1] != `{"b":2}` {
		t.Fatalf("unexpected snapshot: %v", snap)
	}
}

func TestMiddleware_Reset(t *testing.T) {
	m := NewMiddleware(5)
	m.Write(`{"x":1}`)
	m.Reset()
	if m.Len() != 0 {
		t.Fatal("expected 0 after reset")
	}
}

func TestInvalidJSONError_Message(t *testing.T) {
	e := &InvalidJSONError{Line: "bad"}
	if e.Error() != "buffer: invalid JSON: bad" {
		t.Fatalf("unexpected error message: %s", e.Error())
	}
}
