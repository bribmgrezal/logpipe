package batch

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func writeTemp(t *testing.T, v any) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.json")
	if err != nil {
		t.Fatal(err)
	}
	if err := json.NewEncoder(f).Encode(v); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, map[string]any{"size": 50, "interval": "2s"})
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Size != 50 || cfg.Interval != "2s" {
		t.Fatalf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/no/such/file.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "*.json")
	f.WriteString("bad")
	f.Close()
	_, err := LoadConfig(f.Name())
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	_, err := NewFromConfig(nil, func(_ []map[string]any) {})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestNewFromConfig_Valid(t *testing.T) {
	cfg := &Config{Size: 20, Interval: "1s"}
	b, err := NewFromConfig(cfg, func(_ []map[string]any) {})
	if err != nil {
		t.Fatal(err)
	}
	defer b.Stop()
	if b.size != 20 || b.interval != time.Second {
		t.Fatalf("unexpected batcher state")
	}
}

func TestNewFromConfig_InvalidInterval(t *testing.T) {
	cfg := &Config{Size: 10, Interval: "bad"}
	_, err := NewFromConfig(cfg, func(_ []map[string]any) {})
	if err == nil {
		t.Fatal("expected error")
	}
}
