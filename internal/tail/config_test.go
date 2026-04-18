package tail

import (
	"encoding/json"
	"os"
	"testing"
)

func writeTemp(t *testing.T, v any) string {
	t.Helper()
	data, _ := json.Marshal(v)
	f, err := os.CreateTemp("", "tail-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Write(data)
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTemp(t, map[string]any{"path": "/var/log/app.log", "poll_ms": 200})
	defer os.Remove(p)
	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Path != "/var/log/app.log" {
		t.Errorf("path mismatch: %s", cfg.Path)
	}
	if cfg.PollMs != 200 {
		t.Errorf("poll_ms mismatch: %d", cfg.PollMs)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/no/such/file.json")
	if err == nil {
		t.Error("expected error")
	}
}

func TestLoadConfig_MissingPath(t *testing.T) {
	p := writeTemp(t, map[string]any{"poll_ms": 100})
	defer os.Remove(p)
	_, err := LoadConfig(p)
	if err == nil {
		t.Error("expected error for missing path")
	}
}

func TestNewFromConfig(t *testing.T) {
	cfg := &Config{Path: "/tmp/x.log", PollMs: 300}
	tlr := NewFromConfig(cfg)
	if tlr.path != "/tmp/x.log" {
		t.Errorf("unexpected path: %s", tlr.path)
	}
}
