package truncate

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoRules(t *testing.T) {
	tr := New(nil)
	in := []byte(`{"msg":"hello world"}`)
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(in) {
		t.Errorf("expected unchanged, got %s", out)
	}
}

func TestApply_TruncatesLongField(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 10}})
	in := []byte(`{"msg":"this is a very long message"}`)
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	val := m["msg"].(string)
	if len(val) != 10 {
		t.Errorf("expected len 10, got %d: %s", len(val), val)
	}
	if val[len(val)-3:] != "..." {
		t.Errorf("expected suffix '...', got %s", val)
	}
}

func TestApply_ShortFieldUnchanged(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 50}})
	in := []byte(`{"msg":"short"}`)
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["msg"] != "short" {
		t.Errorf("unexpected value: %v", m["msg"])
	}
}

func TestApply_CustomSuffix(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 8, Suffix: "~"}})
	in := []byte(`{"msg":"hello world"}`)
	out, err := tr.Apply(in)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	val := m["msg"].(string)
	if val != "hello wo~" {
		t.Errorf("unexpected value: %s", val)
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	tr := New([]Rule{{Field: "msg", MaxLen: 5}})
	_, err := tr.Apply([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error")
	}
}
