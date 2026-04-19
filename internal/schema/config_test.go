package schema

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "schema-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, `{"rules":[{"field":"level","required":true,"type":"string"}]}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Field != "level" {
		t.Fatal("unexpected config")
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	if _, err := LoadConfig("/no/such/file.json"); err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{bad json}`)
	if _, err := LoadConfig(path); err == nil {
		t.Fatal("expected parse error")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	v := NewFromConfig(nil)
	if v == nil {
		t.Fatal("expected non-nil validator")
	}
	if err := v.Validate([]byte(`{"x":1}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
