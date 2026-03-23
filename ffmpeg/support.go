package ffmpeg

import (
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// AVRational represents a rational number (num/den), re-exported from raw.
// This is a value type used throughout FFmpeg for time bases, frame rates, etc.
type AVRational = raw.AVRational

// cString converts a Go string to a null-terminated C string (*byte).
// The returned pointer is only valid for the duration of the current
// function call — do not store it.
func cString(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	// b[len(s)] is already 0
	return &b[0]
}

// goString converts a null-terminated C string to a Go string.
func goString(p *byte) string {
	if p == nil {
		return ""
	}
	ptr := unsafe.Pointer(p)
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, length)) == 0 {
			break
		}
		length++
	}
	return string(unsafe.Slice(p, length))
}

// avError wraps an FFmpeg AVERROR code into a Go error.
type avError int32

// AVERROR codes — computed from FFERRTAG macro: -(chr0 | chr1<<8 | chr2<<16 | chr3<<24)
const (
	averrorEOF         = -('E' | 'O'<<8 | 'F'<<16 | ' '<<24)
	averrorEAGAIN      = -11 // POSIX EAGAIN
	averrorInvalidData = -('I' | 'N'<<8 | 'D'<<16 | 'A'<<24)
)

// Sentinel errors for common AVERROR codes.
var (
	ErrEOF         = avError(averrorEOF)
	ErrEAGAIN      = avError(averrorEAGAIN)
	ErrInvalidData = avError(averrorInvalidData)
)

func (e avError) Error() string {
	buf := make([]byte, 64)
	ret := raw.AVStrerror(int32(e), &buf[0], uintptr(len(buf)))
	if ret < 0 {
		return "unknown error"
	}
	for i, b := range buf {
		if b == 0 {
			return string(buf[:i])
		}
	}
	return string(buf)
}

// Code returns the raw AVERROR integer code.
func (e avError) Code() int32 { return int32(e) }

// avErr converts an FFmpeg int32 return code to a Go error.
// Returns nil if ret >= 0 (success).
func avErr(ret int32) error {
	if ret >= 0 {
		return nil
	}
	return avError(ret)
}
