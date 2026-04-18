package transform

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
	tr := New(nil)
	input := `{"level":"info","msg":"hello"}`
	out, err := tr.Apply(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != input {
		t.Errorf("expected unchanged output, got %s", out)
	}
}

func TestApply_RenameField(t *testing.T) {
	tr := New([]Rule{{Action: "rename", Field: "msg", Value: "message"}})
	out, err := tr.Apply(`{"msg":"hello"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["msg"]; ok {
		t.Error("old field 'msg' should not exist")
	}
	if m["message"] != "hello" {
		t.Errorf("expected message=hello, got %v", m["message"])
	}
}

func TestApply_DeleteField(t *testing.T) {
	tr := New([]Rule{{Action: "delete", Field: "secret"}})
	out, err := tr.Apply(`{"level":"info","secret":"s3cr3t"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if _, ok := m["secret"]; ok {
		t.Error("field 'secret' should have been deleted")
	}
}

func TestApply_AddField(t *testing.T) {
	tr := New([]Rule{{Action: "add", Field: "env", Value: "production"}})
	out, err := tr.Apply(`{"level":"warn"}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := decode(t, out)
	if m["env"] != "production" {
		t.Errorf("expected env=production, got %v", m["env"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr := New([]Rule{{Action: "delete", Field: "x"}})
	_, err := tr.Apply(`not json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
