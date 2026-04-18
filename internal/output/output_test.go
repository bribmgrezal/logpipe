package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestStdoutWriter_Write(t *testing.T) {
	buf := &bytes.Buffer{}
	w := &StdoutWriter{w: buf}

	if err := w.Write([]byte(`{"level":"info","msg":"hello"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "{\"level\":\"info\",\"msg\":\"hello\"}\n"
	if buf.String() != expected {
		t.Errorf("got %q, want %q", buf.String(), expected)
	}
}

func TestStdoutWriter_Close(t *testing.T) {
	w := NewStdoutWriter()
	if err := w.Close(); err != nil {
		t.Errorf("Close() returned error: %v", err)
	}
}

func TestFileWriter_WriteAndClose(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.log")

	fw, err := NewFileWriter(path)
	if err != nil {
		t.Fatalf("NewFileWriter error: %v", err)
	}

	line := []byte(`{"level":"error","msg":"oops"}`)
	if err := fw.Write(line); err != nil {
		t.Fatalf("Write error: %v", err)
	}
	if err := fw.Close(); err != nil {
		t.Fatalf("Close error: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	expected := string(line) + "\n"
	if string(data) != expected {
		t.Errorf("got %q, want %q", string(data), expected)
	}
}

func TestNewFileWriter_InvalidPath(t *testing.T) {
	_, err := NewFileWriter("/nonexistent/dir/test.log")
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
}
