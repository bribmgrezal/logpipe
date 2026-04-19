package tail

import (
	"context"
	"os"
	"testing"
	"time"
)

func writeLine(t *testing.T, f *os.File, line string) {
	t.Helper()
	_, err := f.WriteString(line + "\n")
	if err != nil {
		t.Fatalf("write: %v", err)
	}
}

func TestTailer_EmitsNewLines(t *testing.T) {
	f, err := os.CreateTemp("", "tail-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(f.Name())

	tlr := New(f.Name(), 50*time.Millisecond)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go tlr.Run(ctx) //nolint

	time.Sleep(80 * time.Millisecond)
	writeLine(t, f, `{"level":"info"}`)
	writeLine(t, f, `{"level":"error"}`)

	got := []string{}
	timeout := time.After(1 * time.Second)
	for len(got) < 2 {
		select {
		case line := <-tlr.Lines():
			got = append(got, line)
		case <-timeout:
			t.Fatalf("timeout waiting for lines, got %d", len(got))
		}
	}

	if got[0] != `{"level":"info"}` {
		t.Errorf("unexpected line: %s", got[0])
	}
	if got[1] != `{"level":"error"}` {
		t.Errorf("unexpected line: %s", got[1])
	}
}

func TestTailer_InvalidPath(t *testing.T) {
	tlr := New("/nonexistent/path/file.log", 50*time.Millisecond)
	ctx := context.Background()
	err := tlr.Run(ctx)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestNew_DefaultPollInterval(t *testing.T) {
	tlr := New("/tmp/x", 0)
	if tlr.pollInterval != 500*time.Millisecond {
		t.Errorf("expected default 500ms, got %v", tlr.pollInterval)
	}
}
