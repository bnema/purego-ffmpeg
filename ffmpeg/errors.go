package ffmpeg

import "github.com/bnema/purego-ffmpeg/internal/core"

// AvError wraps an FFmpeg AVERROR code into a Go error.
type AvError = core.AvError

// CheckError converts an FFmpeg int32 return code to a Go error.
var CheckError = core.CheckError

// Sentinel FFmpeg error codes.
var (
	ErrEOF         = core.ErrEOF
	ErrEAGAIN      = core.ErrEAGAIN
	ErrInvalidData = core.ErrInvalidData
)
