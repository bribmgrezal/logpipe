package flatten

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_FlatObject(t *testing.T) {
	f := New(".")
	out, err := f.Apply(`{"a":1,"b":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["a"] != float64(1) || m["b"] != "hello" {
		t.Errorf("unexpected output: %v", m)
	}
}

func TestApply_NestedObject(t *testing.T) {
	f := New(".")
	out, err := f.Apply(`{"a":{"b":{"c":42}}}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["a.b.c"] != float64(42) {
		t.Errorf("expected a.b.c=42, got %v", m)
	}
}

func TestApply_CustomSeparator(t *testing.T) {
	f := New("_")
	out, err := f.Apply(`{"x":{"y":"val"}}`)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["x_y"] != "val" {
		t.Errorf("expected x_y=val, got %v", m)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	f := New(".")
	_, err := f.Apply(`not json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNew_DefaultSeparator(t *testing.T) {
	f := New("")
	if f.separator != "." {
		t.Errorf("expected default separator '.', got %q", f.separator)
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	f := NewFromConfig(nil)
	if f.separator != "." {
		t.Errorf("expected '.', got %q", f.separator)
	}
}
