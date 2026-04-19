package parse

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestNew_DefaultFormat(t *testing.T) {
	p, err := New("")
	if err != nil || p.format != "json" {
		t.Fatalf("expected json format, got err=%v fmt=%s", err, p.format)
	}
}

func TestNew_UnsupportedFormat(t *testing.T) {
	_, err := New("csv")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestApply_ValidJSON(t *testing.T) {
	p, _ := New("json")
	out, err := p.Apply([]byte(`{"level":"info","msg":"ok"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := New("json")
	_, err := p.Apply([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestApply_Logfmt(t *testing.T) {
	p, _ := New("logfmt")
	out, err := p.Apply([]byte(`level=info msg="hello world" ok`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["level"] != "info" {
		t.Errorf("expected level=info, got %v", m["level"])
	}
	if m["ok"] != true {
		t.Errorf("expected ok=true, got %v", m["ok"])
	}
}

func TestApply_LogfmtQuotedValue(t *testing.T) {
	p, _ := New("logfmt")
	out, err := p.Apply([]byte(`msg="hello world"`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["msg"] != "hello world" {
		t.Errorf("expected msg='hello world', got %v", m["msg"])
	}
}
