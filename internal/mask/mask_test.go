package mask

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, s string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	m, _ := New(nil)
	out, err := m.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatal(err)
	}
	if decode(t, out)["msg"] != "hello" {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestApply_MasksField(t *testing.T) {
	m, _ := New([]Rule{{Field: "email", Pattern: `[^@]+@[^@]+`, Mask: "***"}})
	out, err := m.Apply(`{"email":"user@example.com"}`)
	if err != nil {
		t.Fatal(err)
	}
	if decode(t, out)["email"] == "user@example.com" {
		t.Errorf("email should be masked, got: %s", out)
	}
}

func TestApply_CustomMask(t *testing.T) {
	m, _ := New([]Rule{{Field: "token", Pattern: `\w+`, Mask: "[REDACTED]"}})
	out, err := m.Apply(`{"token":"abc123"}`)
	if err != nil {
		t.Fatal(err)
	}
	if decode(t, out)["token"] != "[REDACTED]" {
		t.Errorf("unexpected token: %s", out)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	m, _ := New(nil)
	_, err := m.Apply(`not-json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestApply_DefaultMask(t *testing.T) {
	m, _ := New([]Rule{{Field: "secret", Pattern: `.+`}})
	out, err := m.Apply(`{"secret":"mysecret"}`)
	if err != nil {
		t.Fatal(err)
	}
	if decode(t, out)["secret"] != "***" {
		t.Errorf("expected default mask, got: %s", out)
	}
}

func TestNew_InvalidPattern(t *testing.T) {
	_, err := New([]Rule{{Field: "f", Pattern: `[invalid`}})
	if err == nil {
		t.Error("expected error for invalid regex")
	}
}
