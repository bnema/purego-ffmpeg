package ffmpeg

import (
	"runtime"
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// AvioOpen creates and initializes an AVIOContext for accessing the
// resource indicated by url.
func AvioOpen(s unsafe.Pointer, uRL string, flags int32) int32 {
	uRLC, uRLBuf := cString(uRL)
	defer runtime.KeepAlive(uRLBuf)
	return raw.AvioOpen(s, uRLC, flags)
}

// AvioClosep closes the resource accessed by the AVIOContext *s, frees it,
// and sets the pointer pointing to it to NULL.
func AvioClosep(s unsafe.Pointer) int32 {
	return raw.AvioClosep(s)
}

// AvioAllocContext creates a custom I/O context with Go callbacks.
// bufSize is the internal buffer size (a buffer is allocated via av_malloc).
// writeable should be true for write contexts, false for read.
// opaque is passed through to the callbacks.
// readFn, writeFn, seekFn are uintptr values from purego.NewCallback().
// Pass 0 for any unused callback.
// The caller must keep the callback functions alive for the context lifetime.
func AvioAllocContext(bufSize int, writeable bool, opaque unsafe.Pointer,
	readFn, writeFn, seekFn uintptr) unsafe.Pointer {

	buf := raw.AVMalloc(uintptr(bufSize))
	if buf == nil {
		return nil
	}

	writeFlag := int32(0)
	if writeable {
		writeFlag = 1
	}

	return raw.AvioAllocContext(buf, int32(bufSize), writeFlag, opaque,
		readFn, writeFn, seekFn)
}

// AvioContextFree frees an AVIOContext and its internal buffer.
func AvioContextFree(ctx unsafe.Pointer) {
	raw.AvioContextFree(ctx)
}
