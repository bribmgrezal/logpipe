package extract

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	e := New(nil)
	input := `{"msg":"hello world"}`
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_SplitsField(t *testing.T) {
	e := New([]Rule{
		{Field: "host_port", Delimiter: ":", Targets: []string{"host", "port"}},
	})
	out, err := e.Apply(`{"host_port":"localhost:8080"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["host"] != "localhost" {
		t.Errorf("expected host=localhost, got %v", m["host"])
	}
	if m["port"] != "8080" {
		t.Errorf("expected port=8080, got %v", m["port"])
	}
}

func TestApply_DefaultSpaceDelimiter(t *testing.T) {
	e := New([]Rule{
		{Field: "name", Targets: []string{"first", "last"}},
	})
	out, err := e.Apply(`{"name":"Jane Doe"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["first"] != "Jane" || m["last"] != "Doe" {
		t.Errorf("unexpected split result: %v", m)
	}
}

func TestApply_MissingField(t *testing.T) {
	e := New([]Rule{
		{Field: "missing", Delimiter: ",", Targets: []string{"a", "b"}},
	})
	input := `{"level":"info"}`
	out, err := e.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["a"]; ok {
		t.Error("expected target field 'a' to be absent")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	e := New([]Rule{
		{Field: "f", Delimiter: "-", Targets: []string{"x"}},
	})
	out, err := e.Apply(`not-json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if out != `not-json` {
		t.Errorf("expected original line returned, got %s", out)
	}
}

func TestApply_FewerPartsThanTargets(t *testing.T) {
	e := New([]Rule{
		{Field: "val", Delimiter: ",", Targets: []string{"a", "b", "c"}},
	})
	out, err := e.Apply(`{"val":"x,y"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["a"] != "x" || m["b"] != "y" {
		t.Errorf("unexpected values: %v", m)
	}
	if _, ok := m["c"]; ok {
		t.Error("expected 'c' to be absent when not enough parts")
	}
}
