package lookup

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
	l := New(nil)
	input := []byte(`{"level":"info"}`)
	out, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output")
	}
}

func TestApply_LooksUpField(t *testing.T) {
	l := New([]Rule{
		{Field: "env", Table: map[string]string{"prod": "production", "dev": "development"}},
	})
	input := []byte(`{"env":"prod","msg":"hello"}`)
	out, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env_lookup"] != "production" {
		t.Errorf("expected env_lookup=production, got %v", m["env_lookup"])
	}
}

func TestApply_CustomTarget(t *testing.T) {
	l := New([]Rule{
		{Field: "code", Table: map[string]string{"200": "OK"}, Target: "status_text"},
	})
	input := []byte(`{"code":"200"}`)
	out, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["status_text"] != "OK" {
		t.Errorf("expected status_text=OK, got %v", m["status_text"])
	}
}

func TestApply_MissingField(t *testing.T) {
	l := New([]Rule{
		{Field: "region", Table: map[string]string{"us": "United States"}},
	})
	input := []byte(`{"level":"warn"}`)
	out, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["region_lookup"]; ok {
		t.Errorf("expected no region_lookup key")
	}
}

func TestApply_NoTableMatch(t *testing.T) {
	l := New([]Rule{
		{Field: "env", Table: map[string]string{"prod": "production"}},
	})
	input := []byte(`{"env":"staging"}`)
	out, err := l.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["env_lookup"]; ok {
		t.Errorf("expected no env_lookup for unmatched value")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	l := New([]Rule{{Field: "env", Table: map[string]string{"x": "y"}}})
	_, err := l.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	l, err := NewFromConfig(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Lookup")
	}
}
