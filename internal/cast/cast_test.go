package cast

import (
	"encoding/json"
	"testing"
)

func decode(t *testing.T, line string) map[string]interface{} {
	t.Helper()
	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(line), &obj); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return obj
}

func TestApply_NoRules(t *testing.T) {
	c := New(nil)
	out, err := c.Apply(`{"count":"42"}`)
	if err != nil {
		t.Fatal(err)
	}
	if out != `{"count":"42"}` {
		t.Errorf("unexpected output: %s", out)
	}
}

func TestApply_CastToInt(t *testing.T) {
	c := New([]Rule{{Field: "count", Target: "int"}})
	out, err := c.Apply(`{"count":"99"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["count"] != float64(99) {
		t.Errorf("expected 99, got %v", obj["count"])
	}
}

func TestApply_CastToFloat(t *testing.T) {
	c := New([]Rule{{Field: "ratio", Target: "float"}})
	out, err := c.Apply(`{"ratio":"3.14"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["ratio"] != 3.14 {
		t.Errorf("expected 3.14, got %v", obj["ratio"])
	}
}

func TestApply_CastToString(t *testing.T) {
	c := New([]Rule{{Field: "code", Target: "string"}})
	out, err := c.Apply(`{"code":200}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["code"] != "200" {
		t.Errorf("expected \"200\", got %v", obj["code"])
	}
}

func TestApply_CastToBool(t *testing.T) {
	c := New([]Rule{{Field: "active", Target: "bool"}})
	out, err := c.Apply(`{"active":"true"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if obj["active"] != true {
		t.Errorf("expected true, got %v", obj["active"])
	}
}

func TestApply_MissingField(t *testing.T) {
	c := New([]Rule{{Field: "missing", Target: "int"}})
	out, err := c.Apply(`{"other":"val"}`)
	if err != nil {
		t.Fatal(err)
	}
	obj := decode(t, out)
	if _, ok := obj["missing"]; ok {
		t.Error("missing field should not be added")
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	c := New([]Rule{{Field: "x", Target: "int"}})
	_, err := c.Apply(`not-json`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}
