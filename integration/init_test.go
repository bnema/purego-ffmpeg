//go:build ffmpeg_integration

package integration

import (
	"testing"

	"github.com/bnema/purego-ffmpeg/ffmpeg"
)

func TestInit(t *testing.T) {
	if err := ffmpeg.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}
}

func TestAvError(t *testing.T) {
	if err := ffmpeg.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Test that AVERROR error messages work after Init injects Strerror
	err := ffmpeg.ErrEOF
	msg := err.Error()
	if msg == "" || msg == "unknown error" {
		t.Errorf("ErrEOF.Error() = %q, expected meaningful message", msg)
	}
	t.Logf("ErrEOF message: %q", msg)
}
