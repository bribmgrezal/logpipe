package select_

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

func TestApply_NoFields_PassThrough(t *testing.T) {
	s := New(nil)
	input := []byte(`{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_SelectsFields(t *testing.T) {
	s := New([]string{"level", "msg"})
	input := []byte(`{"level":"info","msg":"hello","ts":"2024-01-01"}`)
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["ts"]; ok {
		t.Error("expected 'ts' to be removed")
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
	if m["msg"] != "hello" {
		t.Errorf("expected msg=hello, got %v", m["msg"])
	}
}

func TestApply_MissingFieldIgnored(t *testing.T) {
	s := New([]string{"level", "nonexistent"})
	input := []byte(`{"level":"warn","msg":"test"}`)
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["nonexistent"]; ok {
		t.Error("expected missing field to be absent from output")
	}
	if m["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", m["level"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s := New([]string{"level"})
	_, err := s.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_EmptyFields_AllRetained(t *testing.T) {
	s := New([]string{})
	input := []byte(`{"a":1,"b":2}`)
	out, err := s.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}
