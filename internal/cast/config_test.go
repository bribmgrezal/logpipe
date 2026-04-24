package cast

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "cast.json")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTemp(t, `{"rules":[{"field":"count","target":"int"}]}`)
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Field != "count" || cfg.Rules[0].Target != "int" {
		t.Errorf("unexpected rule: %+v", cfg.Rules[0])
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/cast.json")
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	p := writeTemp(t, `{invalid}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	c := NewFromConfig(nil)
	if c == nil {
		t.Error("expected non-nil Caster")
	}
	out, err := c.Apply(`{"x":1}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"x":1}` {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestNewFromConfig_WithRules(t *testing.T) {
	cfg := &Config{Rules: []Rule{{Field: "val", Target: "string"}}}
	c := NewFromConfig(cfg)
	if c == nil {
		t.Error("expected non-nil Caster")
	}
	out, err := c.Apply(`{"val":42}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["val"] != "42" {
		t.Errorf("expected \"42\", got %v", obj["val"])
	}
}
