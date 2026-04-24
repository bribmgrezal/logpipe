package drop

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "drop-cfg-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, `{"rules":[{"field":"level","operator":"eq","value":"debug"}]}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Field != "level" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{invalid}`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	d := NewFromConfig(nil)
	if d == nil {
		t.Fatal("expected non-nil Dropper")
	}
	out, err := d.Apply(`{"level":"info"}`)
	if err != nil || out == "" {
		t.Fatalf("expected pass-through, got %q %v", out, err)
	}
}

func TestNewFromConfig_WithRules(t *testing.T) {
	cfg := &Config{
		Rules: []Rule{{Field: "level", Operator: "eq", Value: "debug"}},
	}
	d := NewFromConfig(cfg)
	out, err := d.Apply(`{"level":"debug"}`)
	if err != nil || out != "" {
		t.Fatalf("expected drop, got %q %v", out, err)
	}
}
