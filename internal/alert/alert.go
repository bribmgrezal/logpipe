package alert

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Rule defines a condition and destination for alerting.
type Rule struct {
	Field    string `json:"field"`
	Contains string `json:"contains"`
	Label    string `json:"label"`
}

// Alerter checks log lines against rules and writes alerts.
type Alerter struct {
	rules  []Rule
	writer io.Writer
}

// New creates an Alerter with the given rules and output writer.
func New(rules []Rule, w io.Writer) *Alerter {
	return &Alerter{rules: rules, writer: w}
}

// Check evaluates a JSON log line against all rules.
// If a rule matches, an alert line is written to the writer.
func (a *Alerter) Check(line []byte) error {
	if len(a.rules) == 0 {
		return nil
	}
	var obj map[string]interface{}
	if err := json.Unmarshal(line, &obj); err != nil {
		return fmt.Errorf("alert: invalid json: %w", err)
	}
	for _, r := range a.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		s, ok := val.(string)
		if !ok {
			continue
		}
		if strings.Contains(s, r.Contains) {
			label := r.Label
			if label == "" {
				label = "ALERT"
			}
			fmt.Fprintf(a.writer, "[%s] field=%q contains=%q line=%s\n", label, r.Field, r.Contains, line)
		}
	}
	return nil
}
