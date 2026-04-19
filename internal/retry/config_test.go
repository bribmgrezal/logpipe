package retry

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = os.WriteFile(f.Name(), []byte(content), 0644)
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, `{"max_attempts":5,"delay_ms":100}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.MaxAttempts != 5 {
		t.Fatalf("expected 5, got %d", cfg.MaxAttempts)
	}
	if cfg.DelayMs != 100 {
		t.Fatalf("expected 100, got %d", cfg.DelayMs)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `not json`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	r := NewFromConfig(nil, func(b []byte) error { return nil })
	if r.maxAttempts != 1 {
		t.Fatalf("expected 1, got %d", r.maxAttempts)
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := &Config{MaxAttempts: 4, DelayMs: 50}
	r := NewFromConfig(cfg, func(b []byte) error { return nil })
	if r.maxAttempts != 4 {
		t.Fatalf("expected 4, got %d", r.maxAttempts)
	}
}
