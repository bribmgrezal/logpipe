package rename

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
	p := writeTemp(t, `{"rules":[{"from":"msg","to":"message"}]}`)
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].From != "msg" || cfg.Rules[0].To != "message" {
		t.Errorf("unexpected rule: %+v", cfg.Rules[0])
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path/config.json")
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
	r := NewFromConfig(nil)
	if r == nil {
		t.Fatal("expected non-nil Renamer")
	}
	input := []byte(`{"key":"value"}`)
	out, err := r.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != string(input) {
		t.Errorf("expected unchanged output")
	}
}
