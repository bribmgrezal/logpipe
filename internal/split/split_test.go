package split

import (
	"encoding/json"
	"testing"
)

func encode(t *testing.T, m map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestApply_NoRules(t *testing.T) {
	s := New(nil, "default")
	out, err := s.Apply(encode(t, map[string]interface{}{"level": "info"}))
	if err != nil {
		t.Fatal(err)
	}
	if out != "default" {
		t.Fatalf("expected default, got %s", out)
	}
}

func TestApply_MatchingRule(t *testing.T) {
	s := New([]Rule{{Field: "level", Value: "error", Output: "errors"}}, "default")
	out, err := s.Apply(encode(t, map[string]interface{}{"level": "error", "msg": "oops"}))
	if err != nil {
		t.Fatal(err)
	}
	if out != "errors" {
		t.Fatalf("expected errors, got %s", out)
	}
}

func TestApply_FallbackOnNoMatch(t *testing.T) {
	s := New([]Rule{{Field: "level", Value: "error", Output: "errors"}}, "default")
	out, err := s.Apply(encode(t, map[string]interface{}{"level": "info"}))
	if err != nil {
		t.Fatal(err)
	}
	if out != "default" {
		t.Fatalf("expected default, got %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	s := New(nil, "default")
	_, err := s.Apply([]byte("not-json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	s := NewFromConfig(nil)
	if s == nil {
		t.Fatal("expected non-nil splitter")
	}
}
