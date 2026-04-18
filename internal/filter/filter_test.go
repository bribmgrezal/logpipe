package filter

import "testing"

func TestMatch_NoRules(t *testing.T) {
	f := New(nil)
	if !f.Match(`{"level":"info","msg":"hello"}`) {
		t.Fatal("expected match with no rules")
	}
}

func TestMatch_EqRule(t *testing.T) {
	f := New([]Rule{{Field: "level", Operator: "eq", Value: "error"}})

	if f.Match(`{"level":"info","msg":"ok"}`) {
		t.Fatal("expected no match for level=info")
	}
	if !f.Match(`{"level":"error","msg":"boom"}`) {
		t.Fatal("expected match for level=error")
	}
}

func TestMatch_ContainsRule(t *testing.T) {
	f := New([]Rule{{Field: "msg", Operator: "contains", Value: "timeout"}})

	if !f.Match(`{"msg":"connection timeout reached"}`) {
		t.Fatal("expected match")
	}
	if f.Match(`{"msg":"all good"}`) {
		t.Fatal("expected no match")
	}
}

func TestMatch_ExistsRule(t *testing.T) {
	f := New([]Rule{{Field: "trace_id", Operator: "exists"}})

	if !f.Match(`{"trace_id":"abc123","msg":"req"}`) {
		t.Fatal("expected match when field exists")
	}
	if f.Match(`{"msg":"no trace"}`) {
		t.Fatal("expected no match when field absent")
	}
}

func TestMatch_InvalidJSON(t *testing.T) {
	f := New([]Rule{{Field: "level", Operator: "eq", Value: "info"}})
	if f.Match(`not json`) {
		t.Fatal("expected no match for invalid JSON")
	}
}
