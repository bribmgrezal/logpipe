package coalesce_test

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/yourorg/logpipe/internal/coalesce"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "coalesce-config-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	content := `{
		"rules": [
			{"target": "primary", "fields": ["field_a", "field_b", "field_c"]}
		]
	}`
	path := writeTemp(t, content)

	cfg, err := coalesce.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
	if cfg.Rules[0].Target != "primary" {
		t.Errorf("expected target 'primary', got %q", cfg.Rules[0].Target)
	}
	if len(cfg.Rules[0].Fields) != 3 {
		t.Errorf("expected 3 fields, got %d", len(cfg.Rules[0].Fields))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := coalesce.LoadConfig("/nonexistent/path/config.json")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{not valid json`)
	_, err := coalesce.LoadConfig(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON, got nil")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	c := coalesce.NewFromConfig(nil)
	if c == nil {
		t.Fatal("expected non-nil coalescer from nil config")
	}

	// With nil config, Apply should pass through valid JSON unchanged.
	input := `{"level":"info","msg":"hello"}`
	out, err := c.Apply([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["level"] != "info" {
		t.Errorf("expected level 'info', got %v", got["level"])
	}
}

func TestNewFromConfig_WithRules(t *testing.T) {
	content := `{
		"rules": [
			{"target": "resolved_host", "fields": ["hostname", "host", "ip"]}
		]
	}`
	path := writeTemp(t, content)

	cfg, err := coalesce.LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error loading config: %v", err)
	}

	c := coalesce.NewFromConfig(cfg)
	if c == nil {
		t.Fatal("expected non-nil coalescer")
	}

	// 'hostname' is missing, 'host' has a value — should coalesce into 'resolved_host'
	input := `{"host":"web-01","msg":"started"}`
	out, err := c.Apply([]byte(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got map[string]interface{}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatalf("output is not valid JSON: %v", err)
	}
	if got["resolved_host"] != "web-01" {
		t.Errorf("expected resolved_host 'web-01', got %v", got["resolved_host"])
	}
}
