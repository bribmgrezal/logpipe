package coalesce

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
	c := New(nil)
	input := []byte(`{"a":"1"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_FirstNonEmpty(t *testing.T) {
	c := New([]Rule{{Fields: []string{"a", "b", "c"}, Target: "result"}})
	input := []byte(`{"a":"","b":"hello","c":"world"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["result"] != "hello" {
		t.Errorf("expected 'hello', got %v", m["result"])
	}
}

func TestApply_AllEmpty(t *testing.T) {
	c := New([]Rule{{Fields: []string{"a", "b"}, Target: "result"}})
	input := []byte(`{"a":"","b":""}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["result"]; ok {
		t.Errorf("expected no result field when all sources are empty")
	}
}

func TestApply_MissingFields(t *testing.T) {
	c := New([]Rule{{Fields: []string{"x", "y"}, Target: "result"}})
	input := []byte(`{"a":"1"}`)
	out, err := c.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["result"]; ok {
		t.Errorf("expected no result field when source fields are absent")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c := New([]Rule{{Fields: []string{"a"}, Target: "result"}})
	_, err := c.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	_, err := NewFromConfig(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}
