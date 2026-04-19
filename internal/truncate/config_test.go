package truncate

import (
	"os"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "truncate-cfg-*.json")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Remove(f.Name()) })
	f.WriteString(content)
	f.Close()
	return f.Name()
}

func TestLoadConfig_Valid(t *testing.T) {
	path := writeTemp(t, `{"rules":[{"field":"msg","max_len":20}]}`)
	cfg, err := LoadConfig(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Field != "msg" {
		t.Errorf("unexpected config: %+v", cfg)
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	_, err := LoadConfig("/nonexistent/truncate.json")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_InvalidJSON(t *testing.T) {
	path := writeTemp(t, `{bad json}`)
	_, err := LoadConfig(path)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewFromConfig_NilConfig(t *testing.T) {
	tr := NewFromConfig(nil)
	if tr == nil {
		t.Fatal("expected non-nil truncator")
	}
	out, err := tr.Apply([]byte(`{"msg":"hello"}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != `{"msg":"hello"}` {
		t.Errorf("unexpected: %s", out)
	}
}
