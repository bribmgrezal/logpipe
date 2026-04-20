package label

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
	in := []byte(`{"msg":"hello"}`)
	out, err := l.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(in) {
		t.Errorf("expected unchanged output")
	}
}

func TestApply_AddsLabel(t *testing.T) {
	l := New([]Rule{{Key: "env", Value: "prod"}})
	out, err := l.Apply([]byte(`{"msg":"hi"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %v", m["env"])
	}
}

func TestApply_DoesNotOverwrite(t *testing.T) {
	l := New([]Rule{{Key: "env", Value: "prod"}})
	out, err := l.Apply([]byte(`{"env":"dev"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "dev" {
		t.Errorf("expected env=dev (unchanged), got %v", m["env"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	l := New([]Rule{{Key: "env", Value: "prod"}})
	_, err := l.Apply([]byte(`not-json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApply_MultipleLabels(t *testing.T) {
	l := New([]Rule{
		{Key: "env", Value: "staging"},
		{Key: "region", Value: "us-east-1"},
	})
	out, err := l.Apply([]byte(`{"msg":"ok"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "staging" || m["region"] != "us-east-1" {
		t.Errorf("unexpected labels: %v", m)
	}
}
