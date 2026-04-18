package tail

import (
	"bufio"
	"context"
	"io"
	"os"
	"time"
)

// Tailer watches a file and emits new lines as they are appended.
type Tailer struct {
	path     string
	pollInterval time.Duration
	lines    chan string
	err      error
}

// New creates a Tailer for the given file path.
func New(path string, pollInterval time.Duration) *Tailer {
	if pollInterval <= 0 {
		pollInterval = 500 * time.Millisecond
	}
	return &Tailer{
		path:         path,
		pollInterval: pollInterval,
		lines:        make(chan string, 64),
	}
}

// Lines returns the channel of emitted lines.
func (t *Tailer) Lines() <-chan string {
	return t.lines
}

// Run starts tailing the file until ctx is cancelled.
func (t *Tailer) Run(ctx context.Context) error {
	f, err := os.Open(t.path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Seek to end on start.
	if _, err := f.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	ticker := time.NewTicker(t.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			close(t.lines)
			return ctx.Err()
		case <-ticker.C:
			for {
				line, err := reader.ReadString('\n')
				if len(line) > 0 {
					if line[len(line)-1] == '\n' {
						line = line[:len(line)-1]
					}
					t.lines <- line
				}
				if err == io.EOF {
					break
				}
				if err != nil {
					close(t.lines)
					return err
				}
			}
		}
	}
}
