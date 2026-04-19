package retry

import (
	"errors"
	"testing"
	"time"
)

func TestWrite_SuccessFirstAttempt(t *testing.T) {
	calls := 0
	r := New(3, 0, func(b []byte) error {
		calls++
		return nil
	})
	if err := r.Write([]byte(`{"msg":"ok"}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestWrite_RetriesOnFailure(t *testing.T) {
	calls := 0
	r := New(3, 0, func(b []byte) error {
		calls++
		if calls < 3 {
			return errors.New("fail")
		}
		return nil
	})
	if err := r.Write([]byte(`{"x":1}`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestWrite_ExceedsMaxAttempts(t *testing.T) {
	calls := 0
	r := New(2, 0, func(b []byte) error {
		calls++
		return errors.New("always fail")
	})
	if err := r.Write([]byte(`{"a":1}`)); err == nil {
		t.Fatal("expected error")
	}
	if calls != 2 {
		t.Fatalf("expected 2 calls, got %d", calls)
	}
}

func TestWrite_InvalidJSON(t *testing.T) {
	r := New(3, 0, func(b []byte) error { return nil })
	err := r.Write([]byte(`not json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
	var e *InvalidJSONError
	if !errors.As(err, &e) {
		t.Fatalf("expected InvalidJSONError, got %T", err)
	}
}

func TestNew_DefaultMinAttempts(t *testing.T) {
	r := New(0, 0, func(b []byte) error { return nil })
	if r.maxAttempts != 1 {
		t.Fatalf("expected 1, got %d", r.maxAttempts)
	}
}

func TestWrite_DelayBetweenRetries(t *testing.T) {
	start := time.Now()
	r := New(2, 20*time.Millisecond, func(b []byte) error {
		return errors.New("fail")
	})
	_ = r.Write([]byte(`{"t":1}`))
	if time.Since(start) < 20*time.Millisecond {
		t.Fatal("expected delay between retries")
	}
}
