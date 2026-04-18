package aggregate

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "agg-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString(content)
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	p := writeTemp(t, `{"field":"level","interval":"5s"}`)
	c, err := LoadConfig(p)
	if err != nil {
		t.Fatal(err)
	}
	if c.Field != "level" || c.Interval != "5s" {
		t.Fatalf("unexpected config: %+v", c)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_MissingField(t *testing.T) {
	p := writeTemp(t, `{"interval":"5s"}`)
	_, err := LoadConfig(p)
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestNewFromConfig_InvalidInterval(t *testing.T) {
	c := &Config{Field: "level", Interval: "bad"}
	_, err := NewFromConfig(c, func([]byte) error { return nil })
	if err == nil {
		t.Fatal("expected error for bad interval")
	}
}

func TestNewFromConfig_DefaultInterval(t *testing.T) {
	c := &Config{Field: "level"}
	a, err := NewFromConfig(c, func([]byte) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	if a.interval.Seconds() != 10 {
		t.Fatalf("expected 10s default, got %v", a.interval)
	}
}
