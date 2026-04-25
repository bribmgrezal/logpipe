package merge

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
	m := New(nil)
	input := []byte(`{"level":"info","msg":"hello"}`)
	out, err := m.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_MergesFields(t *testing.T) {
	rules := []Rule{
		{Target: "context", Sources: []string{"host", "env"}, Remove: false},
	}
	m := New(rules)
	input := []byte(`{"host":"web-01","env":"prod","msg":"ok"}`)
	out, err := m.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := decode(t, out)
	ctx, ok := record["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected context to be a map, got %T", record["context"])
	}
	if ctx["host"] != "web-01" || ctx["env"] != "prod" {
		t.Errorf("unexpected context contents: %v", ctx)
	}
	// source fields should still exist
	if _, exists := record["host"]; !exists {
		t.Error("expected host to remain in record")
	}
}

func TestApply_RemovesSourceFields(t *testing.T) {
	rules := []Rule{
		{Target: "meta", Sources: []string{"host", "env"}, Remove: true},
	}
	m := New(rules)
	input := []byte(`{"host":"web-01","env":"prod","msg":"ok"}`)
	out, err := m.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := decode(t, out)
	if _, exists := record["host"]; exists {
		t.Error("expected host to be removed")
	}
	if _, exists := record["env"]; exists {
		t.Error("expected env to be removed")
	}
	if record["meta"] == nil {
		t.Error("expected meta field to be set")
	}
}

func TestApply_MissingSourceSkipped(t *testing.T) {
	rules := []Rule{
		{Target: "context", Sources: []string{"host", "missing"}, Remove: false},
	}
	m := New(rules)
	input := []byte(`{"host":"web-01","msg":"ok"}`)
	out, err := m.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := decode(t, out)
	ctx, ok := record["context"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected context map, got %T", record["context"])
	}
	if _, exists := ctx["missing"]; exists {
		t.Error("expected missing field to be absent from merged object")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m := New([]Rule{{Target: "ctx", Sources: []string{"a"}}})
	_, err := m.Apply([]byte(`not-json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
