package drop

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
	d := New(nil)
	out, err := d.Apply(`{"level":"info"}`)
	if err != nil || out == "" {
		t.Fatalf("expected pass-through, got %q %v", out, err)
	}
}

func TestApply_DropOnEq(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	out, err := d.Apply(`{"level":"debug","msg":"test"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "" {
		t.Fatalf("expected line dropped, got %q", out)
	}
}

func TestApply_NoDropOnMismatch(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	line := `{"level":"info","msg":"hello"}`
	out, err := d.Apply(line)
	if err != nil || out != line {
		t.Fatalf("expected pass-through, got %q %v", out, err)
	}
}

func TestApply_DropOnContains(t *testing.T) {
	d := New([]Rule{{Field: "msg", Operator: "contains", Value: "healthcheck"}})
	out, err := d.Apply(`{"msg":"GET /healthcheck 200"}`)
	if err != nil || out != "" {
		t.Fatalf("expected drop, got %q %v", out, err)
	}
}

func TestApply_DropOnExists(t *testing.T) {
	d := New([]Rule{{Field: "internal", Operator: "exists"}})
	out, err := d.Apply(`{"internal":true,"msg":"skip"}`)
	if err != nil || out != "" {
		t.Fatalf("expected drop, got %q %v", out, err)
	}
}

func TestApply_DropOnMissing(t *testing.T) {
	d := New([]Rule{{Field: "request_id", Operator: "missing"}})
	out, err := d.Apply(`{"msg":"no id"}`)
	if err != nil || out != "" {
		t.Fatalf("expected drop, got %q %v", out, err)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	d := New([]Rule{{Field: "level", Operator: "eq", Value: "debug"}})
	_, err := d.Apply(`not-json`)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
