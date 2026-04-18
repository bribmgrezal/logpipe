package alert

import (
	"bytes"
	"strings"
	"testing"
)

func TestCheck_NoRules(t *testing.T) {
	var buf bytes.Buffer
	a := New(nil, &buf)
	if err := a.Check([]byte(`{"level":"error"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}

func TestCheck_MatchingRule(t *testing.T) {
	var buf bytes.Buffer
	rules := []Rule{{Field: "level", Contains: "error", Label: "ERR"}}
	a := New(rules, &buf)
	if err := a.Check([]byte(`{"level":"error","msg":"boom"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[ERR]") {
		t.Fatalf("expected ERR label in output, got: %s", out)
	}
}

func TestCheck_NoMatch(t *testing.T) {
	var buf bytes.Buffer
	rules := []Rule{{Field: "level", Contains: "error"}}
	a := New(rules, &buf)
	if err := a.Check([]byte(`{"level":"info"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() != 0 {
		t.Fatalf("expected no output, got: %s", buf.String())
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	var buf bytes.Buffer
	rules := []Rule{{Field: "level", Contains: "error"}}
	a := New(rules, &buf)
	if err := a.Check([]byte(`not-json`)); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestCheck_DefaultLabel(t *testing.T) {
	var buf bytes.Buffer
	rules := []Rule{{Field: "msg", Contains: "fail"}}
	a := New(rules, &buf)
	_ = a.Check([]byte(`{"msg":"fail hard"}`))
	if !strings.Contains(buf.String(), "[ALERT]") {
		t.Fatalf("expected default ALERT label, got: %s", buf.String())
	}
}
