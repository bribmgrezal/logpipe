package limit

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, data []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	l := New(nil)
	input := []byte(`{"msg":"hello world"}`)
	out, drop, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drop {
		t.Fatal("expected no drop")
	}
	if string(out) != string(input) {
		t.Fatalf("expected unchanged output, got %s", out)
	}
}

func TestApply_TruncatesField(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	input := []byte(`{"msg":"hello world","level":"info"}`)
	out, drop, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if drop {
		t.Fatal("expected no drop")
	}
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Fatalf("expected truncated msg, got %v", m["msg"])
	}
}

func TestApply_ShortFieldUnchanged(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 100}})
	input := []byte(`{"msg":"hi"}`)
	out, drop, err := l.Apply(input)
	if err != nil || drop {
		t.Fatalf("unexpected error or drop: %v %v", err, drop)
	}
	m := decode(t, out)
	if m["msg"] != "hi" {
		t.Fatalf("expected unchanged msg, got %v", m["msg"])
	}
}

func TestApply_DropLine(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 3, Drop: true}})
	input := []byte(`{"msg":"too long value"}`)
	_, drop, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !drop {
		t.Fatal("expected line to be dropped")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	l := New([]Rule{{Field: "msg", MaxLen: 5}})
	_, _, err := l.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_MissingField(t *testing.T) {
	l := New([]Rule{{Field: "missing", MaxLen: 5}})
	input := []byte(`{"msg":"hello world"}`)
	out, drop, err := l.Apply(input)
	if err != nil || drop {
		t.Fatalf("unexpected error or drop: %v %v", err, drop)
	}
	m := decode(t, out)
	if m["msg"] != "hello world" {
		t.Fatalf("expected msg unchanged, got %v", m["msg"])
	}
}
