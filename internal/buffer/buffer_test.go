package buffer

import (
	"testing"
)

func TestNew_DefaultCapacity(t *testing.T) {
	b := New(0)
	if b.cap != 100 {
		t.Fatalf("expected cap 100, got %d", b.cap)
	}
}

func TestPush_And_Snapshot(t *testing.T) {
	b := New(3)
	b.Push("a")
	b.Push("b")
	b.Push("c")
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 items, got %d", len(snap))
	}
	if snap[0] != "a" || snap[1] != "b" || snap[2] != "c" {
		t.Fatalf("unexpected snapshot order: %v", snap)
	}
}

func TestPush_Overflow(t *testing.T) {
	b := New(3)
	b.Push("a")
	b.Push("b")
	b.Push("c")
	b.Push("d") // overwrites "a"
	snap := b.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 items, got %d", len(snap))
	}
	if snap[0] != "b" || snap[1] != "c" || snap[2] != "d" {
		t.Fatalf("unexpected snapshot after overflow: %v", snap)
	}
}

func TestLen(t *testing.T) {
	b := New(5)
	if b.Len() != 0 {
		t.Fatal("expected len 0")
	}
	b.Push("x")
	if b.Len() != 1 {
		t.Fatal("expected len 1")
	}
}

func TestReset(t *testing.T) {
	b := New(5)
	b.Push("x")
	b.Push("y")
	b.Reset()
	if b.Len() != 0 {
		t.Fatal("expected len 0 after reset")
	}
	if b.Snapshot() != nil {
		t.Fatal("expected nil snapshot after reset")
	}
}

func TestSnapshot_Empty(t *testing.T) {
	b := New(10)
	if b.Snapshot() != nil {
		t.Fatal("expected nil for empty buffer")
	}
}
