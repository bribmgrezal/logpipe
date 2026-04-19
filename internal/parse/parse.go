package parse

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Parser converts raw log lines into structured JSON maps.
type Parser struct {
	format string // "json" | "logfmt"
}

// New returns a Parser for the given format.
func New(format string) (*Parser, error) {
	format = strings.ToLower(strings.TrimSpace(format))
	if format == "" {
		format = "json"
	}
	if format != "json" && format != "logfmt" {
		return nil, fmt.Errorf("parse: unsupported format %q", format)
	}
	return &Parser{format: format}, nil
}

// Apply parses line and returns a JSON-encoded map.
func (p *Parser) Apply(line []byte) ([]byte, error) {
	switch p.format {
	case "json":
		return parseJSON(line)
	case "logfmt":
		return parseLogfmt(line)
	}
	return nil, fmt.Errorf("parse: unknown format")
}

func parseJSON(line []byte) ([]byte, error) {
	var m map[string]interface{}
	if err := json.Unmarshal(line, &m); err != nil {
		return nil, fmt.Errorf("parse: invalid JSON: %w", err)
	}
	out, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func parseLogfmt(line []byte) ([]byte, error) {
	m := map[string]interface{}{}
	parts := strings.Fields(string(line))
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			m[kv[0]] = strings.Trim(kv[1], `"`)
		} else {
			m[kv[0]] = true
		}
	}
	out, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return out, nil
}
