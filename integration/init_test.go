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

func TestInitVersion(t *testing.T) {
	if err := ffmpeg.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	ver := ffmpeg.FormatVersion()
	if ver == 0 {
		t.Fatal("FormatVersion() returned 0")
	}
	t.Logf("avformat version: %d.%d.%d", ver>>16, (ver>>8)&0xFF, ver&0xFF)
}

func TestAvError(t *testing.T) {
	if err := ffmpeg.Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Test that AVERROR error messages work
	err := ffmpeg.ErrEOF
	msg := err.Error()
	if msg == "" || msg == "unknown error" {
		t.Errorf("ErrEOF.Error() = %q, expected meaningful message", msg)
	}
	t.Logf("ErrEOF message: %q", msg)
}
