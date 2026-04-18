package router

import (
	"bytes"
	"strings"
	"testing"
)

func TestRoute_MatchesRule(t *testing.T) {
	var target, fallback bytes.Buffer
	r := New(
		[]Rule{{Field: "level", Value: "error", Target: "errors"}},
		map[string]interface{ Write([]byte) (int, error) }{"errors": &target},
		&fallback,
	)
	err := r.Route(`{"level":"error","msg":"oops"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(target.String(), "oops") {
		t.Errorf("expected target to contain log line, got: %q", target.String())
	}
	if fallback.Len() != 0 {
		t.Errorf("expected fallback to be empty")
	}
}

func TestRoute_FallbackOnNoMatch(t *testing.T) {
	var fallback bytes.Buffer
	r := New(
		[]Rule{{Field: "level", Value: "error", Target: "errors"}},
		map[string]interface{ Write([]byte) (int, error) }{},
		&fallback,
	)
	err := r.Route(`{"level":"info","msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(fallback.String(), "hello") {
		t.Errorf("expected fallback to contain log line, got: %q", fallback.String())
	}
}

func TestRoute_InvalidJSON(t *testing.T) {
	r := New(nil, nil, nil)
	err := r.Route(`not json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestRoute_NoFallbackNoMatch(t *testing.T) {
	r := New(
		[]Rule{{Field: "level", Value: "error", Target: "errors"}},
		map[string]interface{ Write([]byte) (int, error) }{},
		nil,
	)
	err := r.Route(`{"level":"debug","msg":"trace"}`)
	if err != nil {
		t.Errorf("expected no error when no fallback, got: %v", err)
	}
}
