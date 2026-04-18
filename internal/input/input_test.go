package input

import (
	"io"
	"strings"
	"testing"
)

// readCloser wraps a strings.Reader to satisfy io.ReadCloser.
type readCloser struct {
	io.Reader
}

func (rc readCloser) Close() error { return nil }

func lines(r *Reader) []string {
	var result []string
	for line := range r.Lines() {
		result = append(result, line)
	}
	return result
}

func TestReader_EmptyInput(t *testing.T) {
	r := NewReaderFrom(readCloser{strings.NewReader("")})
	got := lines(r)
	if len(got) != 0 {
		t.Fatalf("expected no lines, got %d", len(got))
	}
}

func TestReader_SingleLine(t *testing.T) {
	r := NewReaderFrom(readCloser{strings.NewReader(`{"level":"info","msg":"hello"}`)})
	got := lines(r)
	if len(got) != 1 {
		t.Fatalf("expected 1 line, got %d", len(got))
	}
	if got[0] != `{"level":"info","msg":"hello"}` {
		t.Errorf("unexpected line: %s", got[0])
	}
}

func TestReader_MultipleLines(t *testing.T) {
	input := "{\"a\":1}\n{\"b\":2}\n{\"c\":3}"
	r := NewReaderFrom(readCloser{strings.NewReader(input)})
	got := lines(r)
	if len(got) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(got))
	}
}

func TestNewFileReader_InvalidPath(t *testing.T) {
	_, err := NewFileReader("/nonexistent/path/file.log")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
