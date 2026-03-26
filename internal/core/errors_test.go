package core

import (
	"errors"
	"testing"
)

func TestAvError_Error_WithoutStrerror(t *testing.T) {
	e := AvError(-11)
	got := e.Error()
	if got != "avError(-11)" {
		t.Errorf("got %q, want %q", got, "avError(-11)")
	}
}

func TestAvError_Code(t *testing.T) {
	e := AvError(-42)
	if e.Code() != -42 {
		t.Errorf("got %d, want -42", e.Code())
	}
}

func TestCheckError_Success(t *testing.T) {
	if err := CheckError(0); err != nil {
		t.Errorf("got %v, want nil", err)
	}
	if err := CheckError(100); err != nil {
		t.Errorf("got %v, want nil", err)
	}
}

func TestCheckError_Failure(t *testing.T) {
	err := CheckError(-11)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var av AvError
	if !errors.As(err, &av) {
		t.Fatalf("expected AvError, got %T", err)
	}
	if av.Code() != -11 {
		t.Errorf("got code %d, want -11", av.Code())
	}
}

func TestSentinelErrors(t *testing.T) {
	if ErrEOF.Code() != -541478725 {
		t.Errorf("ErrEOF code = %d", ErrEOF.Code())
	}
	if ErrEAGAIN.Code() != -11 {
		t.Errorf("ErrEAGAIN code = %d", ErrEAGAIN.Code())
	}
	if ErrInvalidData.Code() != -1094995529 {
		t.Errorf("ErrInvalidData code = %d", ErrInvalidData.Code())
	}
}
