package flatten

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Flattener flattens nested JSON objects into dot-notation keys.
type Flattener struct {
	separator string
	prefix    string
}

// New returns a Flattener with the given separator (defaults to ".").
func New(separator string) *Flattener {
	if separator == "" {
		separator = "."
	}
	return &Flattener{separator: separator}
}

// Apply reads a JSON line, flattens nested objects, and returns the result.
func (f *Flattener) Apply(line string) (string, error) {
	var input map[string]interface{}
	if err := json.Unmarshal([]byte(line), &input); err != nil {
		return "", fmt.Errorf("flatten: invalid JSON: %w", err)
	}

	output := make(map[string]interface{})
	f.flattenMap("", input, output)

	b, err := json.Marshal(output)
	if err != nil {
		return "", fmt.Errorf("flatten: marshal error: %w", err)
	}
	return string(b), nil
}

func (f *Flattener) flattenMap(prefix string, input map[string]interface{}, output map[string]interface{}) {
	for k, v := range input {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, f.separator)
		}
		switch val := v.(type) {
		case map[string]interface{}:
			f.flattenMap(key, val, output)
		default:
			output[key] = val
		}
	}
}
