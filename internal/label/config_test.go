package label

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "label-cfg-*.json")
	if err != nil {
		t.Fatalf("create temp: %v", err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, `{"rules":[{"key":"env","value":"prod"}]}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Key != "env" {
		t.Errorf("unexpected rules: %v", cfg.Rules)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/label.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{bad json}`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	l := NewFromConfig(nil)
	if l == nil {
		t.Error("expected non-nil labeler")
	}
	out, err := l.Apply([]byte(`{"x":1}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(out) != `{"x":1}` {
		t.Errorf("unexpected output: %s", out)
	}
}
