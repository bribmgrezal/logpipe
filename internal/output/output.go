package output

import (
	"fmt"
	"io"
	"os"
)

// Writer defines the interface for log output destinations.
type Writer interface {
	Write(line []byte) error
	Close() error
}

// StdoutWriter writes log lines to stdout.
type StdoutWriter struct {
	w io.Writer
}

// NewStdoutWriter creates a new StdoutWriter.
func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{w: os.Stdout}
}

func (s *StdoutWriter) Write(line []byte) error {
	_, err := fmt.Fprintf(s.w, "%s\n", line)
	return err
}

func (s *StdoutWriter) Close() error {
	return nil
}

// FileWriter writes log lines to a file.
type FileWriter struct {
	f *os.File
}

// NewFileWriter opens or creates a file for appending log output.
func NewFileWriter(path string) (*FileWriter, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("output: open file %q: %w", path, err)
	}
	return &FileWriter{f: f}, nil
}

func (fw *FileWriter) Write(line []byte) error {
	_, err := fmt.Fprintf(fw.f, "%s\n", line)
	return err
}

func (fw *FileWriter) Close() error {
	return fw.f.Close()
}
