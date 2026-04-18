package input

import (
	"bufio"
	"io"
	"os"
)

// Reader reads log lines from an io.Reader line by line.
type Reader struct {
	scanner *bufio.Scanner
	source  io.ReadCloser
}

// NewStdinReader creates a Reader that reads from standard input.
func NewStdinReader() *Reader {
	return &Reader{
		scanner: bufio.NewScanner(os.Stdin),
		source:  os.Stdin,
	}
}

// NewFileReader creates a Reader that reads from the given file path.
func NewFileReader(path string) (*Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &Reader{
		scanner: bufio.NewScanner(f),
		source:  f,
	}, nil
}

// NewReaderFrom creates a Reader from an arbitrary io.ReadCloser.
func NewReaderFrom(rc io.ReadCloser) *Reader {
	return &Reader{
		scanner: bufio.NewScanner(rc),
		source:  rc,
	}
}

// Lines returns a channel that emits each line read from the source.
// The channel is closed when the source is exhausted or an error occurs.
func (r *Reader) Lines() <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for r.scanner.Scan() {
			ch <- r.scanner.Text()
		}
	}()
	return ch
}

// Close releases the underlying source.
func (r *Reader) Close() error {
	return r.source.Close()
}
