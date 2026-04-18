package throttle

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
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, `{"window_seconds":5}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.WindowSeconds != 5 {
		t.Fatalf("expected 5, got %d", cfg.WindowSeconds)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig(filepath.Join(t.TempDir(), "missing.json"))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `not-json`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_ZeroWindow(t *testing.T) {
	path := writeTemp(t, `{"window_seconds":0}`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	th := NewFromConfig(nil)
	if th == nil {
		t.Fatal("expected non-nil throttler")
	}
}
