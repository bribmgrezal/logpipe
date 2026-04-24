package unpack

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
	u := New(nil)
	line := []byte(`{"msg":"hello"}`)
	out, err := u.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(line) {
		t.Fatalf("expected unchanged line, got %s", out)
	}
}

func TestApply_UnpacksField(t *testing.T) {
	u := New([]Rule{{Field: "meta", Remove: false}})
	line := []byte(`{"level":"info","meta":"{\"host\":\"web-1\",\"region\":\"us-east\"}"}`)
	out, err := u.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["host"] != "web-1" {
		t.Errorf("expected host=web-1, got %v", m["host"])
	}
	if m["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", m["region"])
	}
	// original field preserved
	if _, ok := m["meta"]; !ok {
		t.Error("expected meta field to be preserved")
	}
}

func TestApply_RemovesFieldAfterUnpack(t *testing.T) {
	u := New([]Rule{{Field: "meta", Remove: true}})
	line := []byte(`{"level":"info","meta":"{\"host\":\"web-1\"}"}`)
	out, err := u.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["meta"]; ok {
		t.Error("expected meta field to be removed")
	}
	if m["host"] != "web-1" {
		t.Errorf("expected host=web-1, got %v", m["host"])
	}
}

func TestApply_FieldNotString(t *testing.T) {
	u := New([]Rule{{Field: "meta", Remove: true}})
	line := []byte(`{"level":"info","meta":42}`)
	out, err := u.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	// non-string field should be left untouched
	if m["meta"] == nil {
		t.Error("expected meta to remain")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	u := New([]Rule{{Field: "meta"}})
	_, err := u.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_FieldNotValidJSON(t *testing.T) {
	u := New([]Rule{{Field: "meta", Remove: true}})
	line := []byte(`{"meta":"just a plain string"}`)
	out, err := u.Apply(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	// field should remain unchanged when its value is not JSON
	if m["meta"] != "just a plain string" {
		t.Errorf("expected meta unchanged, got %v", m["meta"])
	}
}
