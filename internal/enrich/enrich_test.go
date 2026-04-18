package enrich

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
	e := New(nil)
	in := []byte(`{"level":"info"}`)
	out, err := e.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(in) {
		t.Errorf("expected unchanged, got %s", out)
	}
}

func TestApply_StaticField(t *testing.T) {
	e := New([]Rule{{Field: "env", Value: "production"}})
	out, err := e.Apply([]byte(`{"level":"info"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["env"] != "production" {
		t.Errorf("expected env=production, got %v", m["env"])
	}
}

func TestApply_CopyField(t *testing.T) {
	e := New([]Rule{{Field: "svc", CopyOf: "service"}})
	out, err := e.Apply([]byte(`{"service":"api"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["svc"] != "api" {
		t.Errorf("expected svc=api, got %v", m["svc"])
	}
}

func TestApply_CopyMissingField(t *testing.T) {
	e := New([]Rule{{Field: "svc", CopyOf: "missing"}})
	out, err := e.Apply([]byte(`{"level":"warn"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if _, ok := m["svc"]; ok {
		t.Error("expected svc to be absent")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	e := New([]Rule{{Field: "env", Value: "x"}})
	_, err := e.Apply([]byte(`not-json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestWrap_EnrichesLine(t *testing.T) {
	e := New([]Rule{{Field: "host", Value: "node1"}})
	var got []byte
	wrapped := e.Wrap(func(line []byte) error {
		got = line
		return nil
	})
	if err := wrapped([]byte(`{"msg":"ok"}`)); err != nil {
		t.Fatal(err)
	}
	m := decode(t, got)
	if m["host"] != "node1" {
		t.Errorf("expected host=node1, got %v", m["host"])
	}
}
