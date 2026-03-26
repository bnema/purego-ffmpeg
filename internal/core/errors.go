// internal/core/errors.go
package core

import "fmt"

// AvError wraps an FFmpeg AVERROR code into a Go error.
type AvError int32

// Strerror is injected by the composition root (ffmpeg/init_gen.go)
// after Init() registers symbols. Before injection, Error() falls back
// to a numeric format.
var Strerror func(errnum int32, buf *byte, bufSize uintptr) int32

func (e AvError) Error() string {
	if Strerror != nil {
		buf := make([]byte, 256)
		Strerror(int32(e), &buf[0], uintptr(len(buf)))
		return GoString(&buf[0])
	}
	return fmt.Sprintf("avError(%d)", int32(e))
}

// Code returns the raw AVERROR integer code.
func (e AvError) Code() int32 { return int32(e) }

// AVERROR codes — computed from FFERRTAG macro.
var (
	ErrEOF         = AvError(-('E' | 'O'<<8 | 'F'<<16 | ' '<<24))
	ErrEAGAIN      = AvError(-11)
	ErrInvalidData = AvError(-('I' | 'N'<<8 | 'D'<<16 | 'A'<<24))
)

// CheckError converts an FFmpeg int32 return code to a Go error.
// Returns nil if ret >= 0 (success).
func CheckError(ret int32) error {
	if ret >= 0 {
		return nil
	}
	return AvError(ret)
}
