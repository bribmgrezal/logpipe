package schema

import (
	"encoding/json"
	"testing"
)

func encode(m map[string]interface{}) []byte {
	b, _ := json.Marshal(m)
	return b
}

func TestValidate_NoRules(t *testing.T) {
	v := New(nil)
	if err := v.Validate(encode(map[string]interface{}{"a": 1})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RequiredFieldPresent(t *testing.T) {
	v := New([]Rule{{Field: "level", Required: true}})
	if err := v.Validate(encode(map[string]interface{}{"level": "info"})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_RequiredFieldMissing(t *testing.T) {
	v := New([]Rule{{Field: "level", Required: true}})
	if err := v.Validate(encode(map[string]interface{}{"msg": "hi"})); err == nil {
		t.Fatal("expected error for missing required field")
	}
}

func TestValidate_TypeString_Pass(t *testing.T) {
	v := New([]Rule{{Field: "level", Type: "string"}})
	if err := v.Validate(encode(map[string]interface{}{"level": "info"})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_TypeString_Fail(t *testing.T) {
	v := New([]Rule{{Field: "level", Type: "string"}})
	if err := v.Validate(encode(map[string]interface{}{"level": 42})); err == nil {
		t.Fatal("expected type error")
	}
}

func TestValidate_TypeNumber_Pass(t *testing.T) {
	v := New([]Rule{{Field: "code", Type: "number"}})
	if err := v.Validate(encode(map[string]interface{}{"code": 200})); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidate_InvalidJSON(t *testing.T) {
	v := New([]Rule{{Field: "x", Required: true}})
	if err := v.Validate([]byte("not-json")); err == nil {
		t.Fatal("expected error for invalid json")
	}
}
