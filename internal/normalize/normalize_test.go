package normalize

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
	n := New(nil)
	input := []byte(`{"level":"INFO"}`)
	out, err := n.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output")
	}
}

func TestApply_Lowercase(t *testing.T) {
	n := New([]Rule{{Field: "level", Op: "lowercase"}})
	out, err := n.Apply([]byte(`{"level":"ERROR"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "error" {
		t.Errorf("expected 'error', got %v", m["level"])
	}
}

func TestApply_Uppercase(t *testing.T) {
	n := New([]Rule{{Field: "env", Op: "uppercase"}})
	out, err := n.Apply([]byte(`{"env":"production"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "PRODUCTION" {
		t.Errorf("expected 'PRODUCTION', got %v", m["env"])
	}
}

func TestApply_Trim(t *testing.T) {
	n := New([]Rule{{Field: "msg", Op: "trim"}})
	out, err := n.Apply([]byte(`{"msg":"  hello  "}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "hello" {
		t.Errorf("expected 'hello', got %v", m["msg"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	n := New([]Rule{{Field: "level", Op: "lowercase"}})
	_, err := n.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_MissingField(t *testing.T) {
	n := New([]Rule{{Field: "missing", Op: "trim"}})
	out, err := n.Apply([]byte(`{"level":"info"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Errorf("expected level unchanged")
	}
}
