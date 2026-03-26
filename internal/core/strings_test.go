package core

import "testing"

func TestCString(t *testing.T) {
	p := CString("hello")
	defer FreeCString(p)
	if p == nil {
		t.Fatal("CString returned nil")
	}
	got := GoString(p)
	if got != "hello" {
		t.Errorf("round-trip: got %q, want %q", got, "hello")
	}
}

func TestCString_Empty(t *testing.T) {
	p := CString("")
	defer FreeCString(p)
	got := GoString(p)
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}

func TestGoString_Nil(t *testing.T) {
	got := GoString(nil)
	if got != "" {
		t.Errorf("got %q, want empty string", got)
	}
}
