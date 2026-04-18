package metrics

import "testing"

func TestCounter_InitialZero(t *testing.T) {
	c := &Counter{}
	s := c.Snapshot()
	if s.Received != 0 || s.Passed != 0 || s.Filtered != 0 || s.Errors != 0 {
		t.Fatal("expected all counters to be zero initially")
	}
}

func TestCounter_Increments(t *testing.T) {
	c := &Counter{}
	c.IncReceived()
	c.IncReceived()
	c.IncPassed()
	c.IncFiltered()
	c.IncErrors()

	s := c.Snapshot()
	if s.Received != 2 {
		t.Errorf("expected Received=2, got %d", s.Received)
	}
	if s.Passed != 1 {
		t.Errorf("expected Passed=1, got %d", s.Passed)
	}
	if s.Filtered != 1 {
		t.Errorf("expected Filtered=1, got %d", s.Filtered)
	}
	if s.Errors != 1 {
		t.Errorf("expected Errors=1, got %d", s.Errors)
	}
}

func TestCounter_Reset(t *testing.T) {
	c := &Counter{}
	c.IncReceived()
	c.IncPassed()
	c.Reset()

	s := c.Snapshot()
	if s.Received != 0 || s.Passed != 0 {
		t.Fatal("expected counters to be zero after reset")
	}
}

func TestGlobal_NotNil(t *testing.T) {
	if Global == nil {
		t.Fatal("Global counter should not be nil")
	}
}
