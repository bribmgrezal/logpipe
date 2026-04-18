package format

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoTemplate_PrettyJSON(t *testing.T) {
	f := New("")
	out, err := f.Apply([]byte(`{"level":"info","msg":"hello"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestApply_Template(t *testing.T) {
	f := New("{level} {msg}")
	out, err := f.Apply([]byte(`{"level":"warn","msg":"oops"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "warn oops" {
		t.Errorf("expected 'warn oops', got %q", string(out))
	}
}

func TestApply_MissingField(t *testing.T) {
	f := New("{level} {missing}")
	out, err := f.Apply([]byte(`{"level":"debug"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "debug " {
		t.Errorf("expected 'debug ', got %q", string(out))
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	f := New("{level}")
	_, err := f.Apply([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
