package rename

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
	r := New(nil)
	input := []byte(`{"level":"info","msg":"hello"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_RenamesField(t *testing.T) {
	r := New([]Rule{{From: "msg", To: "message"}})
	input := []byte(`{"level":"info","msg":"hello"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["msg"]; ok {
		t.Error("old field 'msg' should have been removed")
	}
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
}

func TestApply_MissingSourceField(t *testing.T) {
	r := New([]Rule{{From: "nonexistent", To: "target"}})
	input := []byte(`{"level":"warn"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["target"]; ok {
		t.Error("target field should not exist when source is missing")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r := New([]Rule{{From: "a", To: "b"}})
	_, err := r.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	r := New([]Rule{
		{From: "msg", To: "message"},
		{From: "lvl", To: "level"},
	})
	input := []byte(`{"msg":"hello","lvl":"info"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}
