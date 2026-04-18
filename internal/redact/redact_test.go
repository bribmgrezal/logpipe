package redact

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal([]byte(line), &m); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	r := New(nil)
	input := `{"msg":"hello","password":"secret"}`
	if got := r.Apply(input); got != input {
		t.Errorf("expected unchanged, got %s", got)
	}
}

func TestApply_RedactsField(t *testing.T) {
	r := New([]Rule{{Field: "password"}})
	out := r.Apply(`{"user":"alice","password":"secret"}`)
	m := decode(t, out)
	if m["password"] != "***" {
		t.Errorf("expected ***, got %v", m["password"])
	}
	if m["user"] != "alice" {
		t.Errorf("user should be unchanged")
	}
}

func TestApply_CustomReplace(t *testing.T) {
	r := New([]Rule{{Field: "token", Replace: "[REDACTED]"}})
	out := r.Apply(`{"token":"abc123"}`)
	m := decode(t, out)
	if m["token"] != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %v", m["token"])
	}
}

func TestApply_NestedField(t *testing.T) {
	r := New([]Rule{{Field: "user.password"}})
	out := r.Apply(`{"user":{"name":"bob","password":"hunter2"}}`)
	m := decode(t, out)
	user := m["user"].(map[string]interface{})
	if user["password"] != "***" {
		t.Errorf("expected nested field redacted, got %v", user["password"])
	}
	if user["name"] != "bob" {
		t.Errorf("name should be unchanged")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	r := New([]Rule{{Field: "password"}})
	input := `not-json`
	if got := r.Apply(input); got != input {
		t.Errorf("expected original line returned on invalid json")
	}
}

func TestApply_MissingField(t *testing.T) {
	r := New([]Rule{{Field: "secret"}})
	input := `{"msg":"hello"}`
	out := r.Apply(input)
	m := decode(t, out)
	if _, ok := m["secret"]; ok {
		t.Errorf("field should not be added")
	}
}
