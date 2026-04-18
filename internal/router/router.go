package router

import (
	"encoding/json"
	"fmt"
	"io"
)

// Rule defines a routing rule: if Field matches Value, send to Target output.
type Rule struct {
	Field  string `json:"field"`
	Value  string `json:"value"`
	Target string `json:"target"`
}

// Router routes log lines to named outputs based on rules.
type Router struct {
	rules   []Rule
	outputs map[string]io.Writer
	fallback io.Writer
}

// New creates a Router with the given rules, named outputs, and a fallback writer.
func New(rules []Rule, outputs map[string]io.Writer, fallback io.Writer) *Router {
	return &Router{
		rules:    rules,
		outputs:  outputs,
		fallback: fallback,
	}
}

// Route parses a JSON log line and writes it to the appropriate output.
func (r *Router) Route(line string) error {
	var record map[string]interface{}
	if err := json.Unmarshal([]byte(line), &record); err != nil {
		return fmt.Errorf("router: invalid json: %w", err)
	}

	for _, rule := range r.rules {
		val, ok := record[rule.Field]
		if !ok {
			continue
		}
		if fmt.Sprintf("%v", val) == rule.Value {
			if w, found := r.outputs[rule.Target]; found {
				_, err := fmt.Fprintln(w, line)
				return err
			}
		}
	}

	if r.fallback != nil {
		_, err := fmt.Fprintln(r.fallback, line)
		return err
	}
	return nil
}
