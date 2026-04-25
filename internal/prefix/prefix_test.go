package prefix

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
	p := New(nil)
	input := `{"msg":"hello"}`
	out, err := p.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_PrependsPrefix(t *testing.T) {
	p := New([]Rule{{Field: "env", Prefix: "prod-"}})
	out, err := p.Apply(`{"env":"us-east"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "prod-us-east" {
		t.Errorf("expected prod-us-east, got %v", m["env"])
	}
}

func TestApply_MissingFieldSkipped(t *testing.T) {
	p := New([]Rule{{Field: "service", Prefix: "svc-"}})
	out, err := p.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["service"]; ok {
		t.Error("expected missing field to remain absent")
	}
}

func TestApply_NonStringFieldSkipped(t *testing.T) {
	p := New([]Rule{{Field: "count", Prefix: "n-"}})
	out, err := p.Apply(`{"count":42}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["count"] != float64(42) {
		t.Errorf("expected count unchanged, got %v", m["count"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p := New([]Rule{{Field: "msg", Prefix: "pre-"}})
	_, err := p.Apply(`not-json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApply_MultipleRules(t *testing.T) {
	p := New([]Rule{
		{Field: "host", Prefix: "h-"},
		{Field: "region", Prefix: "r-"},
	})
	out, err := p.Apply(`{"host":"web01","region":"east"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["host"] != "h-web01" {
		t.Errorf("expected h-web01, got %v", m["host"])
	}
	if m["region"] != "r-east" {
		t.Errorf("expected r-east, got %v", m["region"])
	}
}
