package tag

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

func TestApply_NoRules(t *testing.T) {
	tgr := New(nil)
	input := `{"level":"error"}`
	out, err := tgr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected passthrough, got %s", out)
	}
}

func TestApply_AddsTag(t *testing.T) {
	tgr := New([]Rule{{Field: "level", Match: "error", Tag: "critical"}})
	out, err := tgr.Apply(`{"level":"error"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	tags, ok := m["tags"].([]interface{})
	if !ok || len(tags) != 1 || tags[0] != "critical" {
		t.Errorf("expected tags=[critical], got %v", m["tags"])
	}
}

func TestApply_NoMatchSkips(t *testing.T) {
	tgr := New([]Rule{{Field: "level", Match: "error", Tag: "critical"}})
	out, err := tgr.Apply(`{"level":"info"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["tags"]; ok {
		t.Errorf("expected no tags field, got %v", m["tags"])
	}
}

func TestApply_CustomTarget(t *testing.T) {
	tgr := New([]Rule{{Field: "env", Match: "prod", Tag: "production", Target: "labels"}})
	out, err := tgr.Apply(`{"env":"prod"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	labels, ok := m["labels"].([]interface{})
	if !ok || len(labels) != 1 || labels[0] != "production" {
		t.Errorf("expected labels=[production], got %v", m["labels"])
	}
}

func TestApply_MultipleTagsSameTarget(t *testing.T) {
	tgr := New([]Rule{
		{Field: "level", Match: "error", Tag: "alert"},
		{Field: "level", Match: "error", Tag: "urgent"},
	})
	out, err := tgr.Apply(`{"level":"error"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	tags, ok := m["tags"].([]interface{})
	if !ok || len(tags) != 2 {
		t.Errorf("expected 2 tags, got %v", m["tags"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tgr := New([]Rule{{Field: "level", Match: "error", Tag: "critical"}})
	_, err := tgr.Apply(`not-json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
