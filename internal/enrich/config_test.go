package enrich

import (
	"os"
	"testing"
)

func TestLoadConfig_Valid(t *testing.T) {
	f, err := os.CreateTemp("", "enrich-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString(`{"rules":[{"field":"env","value":"staging"}]}`)
	f.Close()

	cfg, err := LoadConfig(f.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Field != "env" || cfg.Rules[0].Value != "staging" {
		t.Errorf("unexpected rule: %+v", cfg.Rules[0])
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "enrich-bad-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString(`not-json`)
	f.Close()

	_, err = LoadConfig(f.Name())
	if err == nil {
		t.Error("expected parse error")
	}
}

func TestNewFromConfig(t *testing.T) {
	cfg := &Config{Rules: []Rule{{Field: "region", Value: "us-east"}}}
	e := NewFromConfig(cfg)
	out, err := e.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", m["region"])
	}
}
