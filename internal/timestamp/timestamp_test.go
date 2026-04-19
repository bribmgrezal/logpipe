package timestamp

import (
	"encoding/json"
	"testing"
	"time"
)

func decode(t *testing.T, b []byte) map[string]interface{} {
	t.Helper()
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("decode: %v", err)
	}
	return m
}

func TestApply_NoField(t *testing.T) {
	p, _ := New("ts", "", "")
	line := []byte(`{"msg":"hello"}`)
	out, err := p.Apply(line)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != string(line) {
		t.Errorf("expected unchanged line")
	}
}

func TestApply_NormalizesRFC3339(t *testing.T) {
	p, _ := New("ts", "", time.RFC3339)
	input := `{"ts":"2024-01-15T10:30:00+05:00","msg":"ok"}`
	out, err := p.Apply([]byte(input))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	got := m["ts"].(string)
	if got != "2024-01-15T05:30:00Z" {
		t.Errorf("unexpected ts: %s", got)
	}
}

func TestApply_CustomFormat(t *testing.T) {
	p, _ := New("time", "2006-01-02", "01/02/2006")
	out, err := p.Apply([]byte(`{"time":"2024-03-21"}`))
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["time"] != "03/21/2024" {
		t.Errorf("got %v", m["time"])
	}
}

func TestApply_InvalidJSON(t *testing.T) {
	p, _ := New("ts", "", "")
	_, err := p.Apply([]byte(`not-json`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestApply_UnparseableTimestamp(t *testing.T) {
	p, _ := New("ts", "", "")
	line := []byte(`{"ts":"not-a-time"}`)
	out, err := p.Apply(line)
	if err != nil {
		t.Fatal(err)
	}
	m := decode(t, out)
	if m["ts"] != "not-a-time" {
		t.Errorf("expected original value preserved")
	}
}

func TestNew_EmptyField(t *testing.T) {
	_, err := New("", "", "")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}
