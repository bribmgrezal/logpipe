package sampler

import (
	"testing"
)

func TestNew_DefaultRate(t *testing.T) {
	s := New(0)
	if s.rate != 1 {
		t.Fatalf("expected rate 1, got %d", s.rate)
	}
}

func TestAllow_RateOne_AllPass(t *testing.T) {
	s := New(1)
	for i := 0; i < 10; i++ {
		if !s.Allow() {
			t.Fatal("expected all messages to pass with rate=1")
		}
	}
}

func TestAllow_RateN_SamplesCorrectly(t *testing.T) {
	s := New(3)
	results := make([]bool, 9)
	for i := range results {
		results[i] = s.Allow()
	}
	// positions 0,3,6 should be true
	for i, v := range results {
		expect := i%3 == 0
		if v != expect {
			t.Errorf("index %d: got %v, want %v", i, v, expect)
		}
	}
}

func TestReset_ResetsCounter(t *testing.T) {
	s := New(2)
	s.Allow()
	s.Allow()
	s.Reset()
	if !s.Allow() {
		t.Fatal("expected first call after reset to pass")
	}
}

func TestWrap_DropsInvalidJSON(t *testing.T) {
	s := New(1)
	called := false
	wrapped := s.Wrap(func(b []byte) error {
		called = true
		return nil
	})
	_ = wrapped([]byte("not-json"))
	if called {
		t.Fatal("expected invalid JSON to be dropped")
	}
}

func TestWrap_PassesValidJSON(t *testing.T) {
	s := New(1)
	called := false
	wrapped := s.Wrap(func(b []byte) error {
		called = true
		return nil
	})
	_ = wrapped([]byte(`{"level":"info"}`))
	if !called {
		t.Fatal("expected valid JSON to be forwarded")
	}
}
