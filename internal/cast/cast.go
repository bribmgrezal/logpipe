package cast

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Rule defines a single field cast operation.
type Rule struct {
	Field  string `json:"field"`
	Target string `json:"target"` // "string", "int", "float", "bool"
}

// Caster applies type casting rules to JSON log lines.
type Caster struct {
	rules []Rule
}

// New returns a new Caster with the given rules.
func New(rules []Rule) *Caster {
	return &Caster{rules: rules}
}

// Apply parses the JSON line, casts configured fields, and returns updated JSON.
func (c *Caster) Apply(line string) (string, error) {
	if len(c.rules) == 0 {
		return line, nil
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		return "", fmt.Errorf("cast: invalid JSON: %w", err)
	}

	for _, r := range c.rules {
		val, ok := obj[r.Field]
		if !ok {
			continue
		}
		casted, err := castValue(val, r.Target)
		if err != nil {
			continue
		}
		obj[r.Field] = casted
	}

	out, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("cast: marshal error: %w", err)
	}
	return string(out), nil
}

func castValue(val interface{}, target string) (interface{}, error) {
	str := fmt.Sprintf("%v", val)
	switch target {
	case "string":
		return str, nil
	case "int":
		return strconv.ParseInt(str, 10, 64)
	case "float":
		return strconv.ParseFloat(str, 64)
	case "bool":
		return strconv.ParseBool(str)
	default:
		return nil, fmt.Errorf("cast: unknown target type %q", target)
	}
}
