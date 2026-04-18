package pipeline

import (
	"bufio"
	"fmt"
	"io"

	"github.com/user/logpipe/internal/filter"
	"github.com/user/logpipe/internal/output"
)

// Pipeline reads log lines from a reader, applies a filter, and writes
// matching lines to the configured output writer.
type Pipeline struct {
	filter *filter.Filter
	writer output.Writer
}

// New creates a new Pipeline with the given filter and output writer.
func New(f *filter.Filter, w output.Writer) *Pipeline {
	return &Pipeline{filter: f, writer: w}
}

// Run reads lines from r until EOF or error, forwarding matching lines.
// It returns the number of lines forwarded and any write error encountered.
func (p *Pipeline) Run(r io.Reader) (int, error) {
	scanner := bufio.NewScanner(r)
	forwarded := 0

	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}

		matched, err := p.filter.Match(line)
		if err != nil {
			// skip non-JSON or unparseable lines
			continue
		}

		if matched {
			if err := p.writer.Write(line); err != nil {
				return forwarded, fmt.Errorf("pipeline: write: %w", err)
			}
			forwarded++
		}
	}

	if err := scanner.Err(); err != nil {
		return forwarded, fmt.Errorf("pipeline: scan: %w", err)
	}

	return forwarded, nil
}
