package ffmpeg

import "testing"

func TestCString(t *testing.T) {
	s := "hello"
	p := cString(s)
	if p == nil {
		t.Fatal("cString returned nil")
	}
	got := goString(p)
	if got != s {
		t.Errorf("goString(cString(%q)) = %q", s, got)
	}
}

func TestCStringEmpty(t *testing.T) {
	p := cString("")
	got := goString(p)
	if got != "" {
		t.Errorf("goString(cString(%q)) = %q", "", got)
	}
}

func TestGoStringNil(t *testing.T) {
	got := goString(nil)
	if got != "" {
		t.Errorf("goString(nil) = %q, want empty", got)
	}
}

func TestAvError(t *testing.T) {
	err := avError(averrorEOF)
	if err.Code() != averrorEOF {
		t.Errorf("Code() = %d, want %d", err.Code(), averrorEOF)
	}
	// Error() can't call av_strerror without Init(), just check it doesn't panic
	// with a nil function (it will return "unknown error" since raw.AVStrerror is nil)
}

func TestAvErr(t *testing.T) {
	if err := avErr(0); err != nil {
		t.Errorf("avErr(0) = %v, want nil", err)
	}
	if err := avErr(42); err != nil {
		t.Errorf("avErr(42) = %v, want nil", err)
	}
	if err := avErr(-1); err == nil {
		t.Error("avErr(-1) = nil, want error")
	}
}
