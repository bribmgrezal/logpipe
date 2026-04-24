package dedupe

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
		t.Fatalf("writeTemp: %v", err)
	}
	return p
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTemp(t, `{"fields":["id","host"],"max_size":500}`)
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Fields) != 2 || cfg.Fields[0] != "id" {
		t.Fatalf("unexpected fields: %v", cfg.Fields)
	}
	if cfg.MaxSize != 500 {
		t.Fatalf("expected max_size 500, got %d", cfg.MaxSize)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	p := writeTemp(t, `{invalid}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	d := NewFromConfig(nil)
	if d == nil {
		t.Fatal("expected non-nil Deduper")
	}
	if d.maxSize != 10000 {
		t.Fatalf("expected default maxSize, got %d", d.maxSize)
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := &Config{Fields: []string{"trace_id"}, MaxSize: 200}
	d := NewFromConfig(cfg)
	if d.maxSize != 200 {
		t.Fatalf("expected maxSize 200, got %d", d.maxSize)
	}
	if len(d.fields) != 1 || d.fields[0] != "trace_id" {
		t.Fatalf("unexpected fields: %v", d.fields)
	}
}
