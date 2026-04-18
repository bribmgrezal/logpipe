package alert

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
)

func writeTemp(t *testing.T, v interface{}) string {
	t.Helper()
	f, err := os.CreateTemp("", "alert-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, Config{Rules: []Rule{{Field: "level", Contains: "error", Label: "E"}}})
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(cfg.Rules))
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/alert.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp("", "bad-*.json")
	f.WriteString("not json")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	_, err := LoadConfig(f.Name())
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestNewFromConfig(t *testing.T) {
	cfg := &Config{Rules: []Rule{{Field: "level", Contains: "warn", Label: "W"}}}
	var buf bytes.Buffer
	a := NewFromConfig(cfg, &buf)
	if a == nil {
		t.Fatal("expected non-nil alerter")
	}
	_ = a.Check([]byte(`{"level":"warn"}`))
	if buf.Len() == 0 {
		t.Fatal("expected alert output")
	}
}
