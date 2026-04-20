package validate

import (
	"encoding/json"
	"testing"
)

func encode(t *testing.T, m map[string]interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestCheck_NoRules(t *testing.T) {
	v, _ := New(nil)
	if err := v.Check(encode(t, map[string]interface{}{"msg": "hello"})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_PatternMatch(t *testing.T) {
	v, _ := New([]Rule{{Field: "level", Pattern: `^(info|warn|error)$`}})
	if err := v.Check(encode(t, map[string]interface{}{"level": "info"})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheck_PatternNoMatch(t *testing.T) {
	v, _ := New([]Rule{{Field: "level", Pattern: `^(info|warn|error)$`}})
	err := v.Check(encode(t, map[string]interface{}{"level": "debug"}))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCheck_MissingField(t *testing.T) {
	v, _ := New([]Rule{{Field: "level", Pattern: `.*`}})
	err := v.Check(encode(t, map[string]interface{}{"msg": "hi"}))
	if err == nil {
		t.Fatal("expected error for missing field")
	}
}

func TestCheck_MinLen(t *testing.T) {
	v, _ := New([]Rule{{Field: "msg", MinLen: 5}})
	err := v.Check(encode(t, map[string]interface{}{"msg": "hi"}))
	if err == nil {
		t.Fatal("expected min_len error")
	}
}

func TestCheck_MaxLen(t *testing.T) {
	v, _ := New([]Rule{{Field: "msg", MaxLen: 3}})
	err := v.Check(encode(t, map[string]interface{}{"msg": "toolong"}))
	if err == nil {
		t.Fatal("expected max_len error")
	}
}

func TestCheck_InvalidJSON(t *testing.T) {
	v, _ := New([]Rule{{Field: "x", Pattern: `.*`}})
	if err := v.Check([]byte("not json")); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]Rule{{Field: "f", Pattern: `[invalid`}})
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}
