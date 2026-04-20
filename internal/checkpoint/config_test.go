package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "config.json")
	_ = os.WriteFile(p, []byte(content), 0644)
	return p
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTemp(t, `{"path": "/tmp/cp.json"}`)
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Path != "/tmp/cp.json" {
		t.Errorf("unexpected path: %s", cfg.Path)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/config.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	p := writeTemp(t, `{bad json}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLoadConfig_MissingPath(t *testing.T) {
	p := writeTemp(t, `{}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Error("expected error for missing path field")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	_, err := NewFromConfig(nil)
	if err == nil {
		t.Error("expected error for nil config")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := &Config{Path: filepath.Join(t.TempDir(), "cp.json")}
	cp, err := NewFromConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp == nil {
		t.Error("expected non-nil checkpoint")
	}
}
