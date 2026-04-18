package format

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "config.json")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTemp(t, `{"template":"{level} {msg}"}`)
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Template != "{level} {msg}" {
		t.Errorf("unexpected template: %q", cfg.Template)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	p := writeTemp(t, `not json`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	f := NewFromConfig(nil)
	if f == nil {
		t.Fatal("expected non-nil formatter")
	}
}

func TestNewFromConfig_WithTemplate(t *testing.T) {
	cfg := &Config{Template: "{level}"}
	f := NewFromConfig(cfg)
	out, err := f.Apply([]byte(`{"level":"error"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != "error" {
		t.Errorf("expected 'error', got %q", string(out))
	}
}
