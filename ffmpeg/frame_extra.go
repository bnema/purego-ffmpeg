package ffmpeg

import (
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// HWFrameGetBuffer allocates a new hardware-backed frame from an
// AVHWFramesContext.
func HWFrameGetBuffer(hwframeCtx unsafe.Pointer, frame unsafe.Pointer, flags int32) int32 {
	return raw.AVHwframeGetBuffer(hwframeCtx, frame, flags)
}
