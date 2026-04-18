package throttle

import (
	"testing"
	"time"
)

func TestNew_DefaultWindow(t *testing.T) {
	th := New(0)
	if th.window != time.Second {
		t.Fatalf("expected 1s default, got %v", th.window)
	}
}

func TestAllow_FirstLineAllowed(t *testing.T) {
	th := New(time.Minute)
	if !th.Allow([]byte(`{"msg":"hello"}`)) {
		t.Fatal("expected first line to be allowed")
	}
}

func TestAllow_DuplicateSuppressed(t *testing.T) {
	th := New(time.Minute)
	line := []byte(`{"msg":"hello"}`)
	th.Allow(line)
	if th.Allow(line) {
		t.Fatal("expected duplicate to be suppressed")
	}
}

func TestAllow_AfterWindowExpires(t *testing.T) {
	th := New(50 * time.Millisecond)
	line := []byte(`{"msg":"hello"}`)
	th.Allow(line)
	time.Sleep(60 * time.Millisecond)
	if !th.Allow(line) {
		t.Fatal("expected line to be allowed after window expires")
	}
}

func TestReset_ClearsState(t *testing.T) {
	th := New(time.Minute)
	line := []byte(`{"msg":"hello"}`)
	th.Allow(line)
	th.Reset()
	if !th.Allow(line) {
		t.Fatal("expected line to be allowed after reset")
	}
}

func TestWrap_InvalidJSONPassthrough(t *testing.T) {
	th := New(time.Minute)
	var got []byte
	next := func(line []byte) error { got = line; return nil }
	wrapped := Wrap(th, next)
	input := []byte("not-json")
	if err := wrapped(input); err != nil {
		t.Fatal(err)
	}
	if string(got) != string(input) {
		t.Fatal("expected invalid JSON to pass through")
	}
}
