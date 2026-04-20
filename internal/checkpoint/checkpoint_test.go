package checkpoint

import (
	"os"
	"path/filepath"
	"testing"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNew_NoExistingFile(t *testing.T) {
	cp, err := New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cp.Load().Offset != 0 {
		t.Errorf("expected zero offset, got %d", cp.Load().Offset)
	}
}

func TestSave_And_Load(t *testing.T) {
	cp, _ := New(tempPath(t))
	if err := cp.Save("app.log", 42); err != nil {
		t.Fatalf("save failed: %v", err)
	}
	s := cp.Load()
	if s.File != "app.log" || s.Offset != 42 {
		t.Errorf("unexpected state: %+v", s)
	}
}

func TestPersistence_AcrossInstances(t *testing.T) {
	p := tempPath(t)
	cp1, _ := New(p)
	_ = cp1.Save("stream.log", 99)

	cp2, err := New(p)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	s := cp2.Load()
	if s.File != "stream.log" || s.Offset != 99 {
		t.Errorf("persisted state mismatch: %+v", s)
	}
}

func TestReset_ClearsState(t *testing.T) {
	p := tempPath(t)
	cp, _ := New(p)
	_ = cp.Save("x.log", 7)
	if err := cp.Reset(); err != nil {
		t.Fatalf("reset failed: %v", err)
	}
	if cp.Load().Offset != 0 {
		t.Errorf("expected zero offset after reset")
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Errorf("expected file to be removed after reset")
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not-json"), 0644)
	_, err := New(p)
	if err == nil {
		t.Error("expected error for invalid JSON checkpoint file")
	}
}
