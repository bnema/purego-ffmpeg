package ffmpeg

import (
	"runtime"
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// MediaType represents an AVMediaType enum value.
type MediaType int32

const (
	MediaTypeAvmediaTypeUnknown    MediaType = -1
	MediaTypeAvmediaTypeVideo      MediaType = 0
	MediaTypeAvmediaTypeAudio      MediaType = 1
	MediaTypeAvmediaTypeData       MediaType = 2
	MediaTypeAvmediaTypeSubtitle   MediaType = 3
	MediaTypeAvmediaTypeAttachment MediaType = 4
	MediaTypeAvmediaTypeNb         MediaType = 5
)

// Convenience aliases matching FFmpeg naming.
const (
	AVMEDIA_TYPE_VIDEO MediaType = 0
	AVMEDIA_TYPE_AUDIO MediaType = 1
)

// AvutilVersion returns the LIBAVUTIL_VERSION_INT constant.
func AvutilVersion() uint32 {
	return raw.AvutilVersion()
}

// AVMalloc allocates a block of size bytes with av_malloc.
func AVMalloc(size uintptr) unsafe.Pointer {
	return raw.AVMalloc(size)
}

// AVFree frees a block allocated with av_malloc.
func AVFree(ptr unsafe.Pointer) {
	raw.AVFree(ptr)
}

// AVOptSet sets a string option on an AVClass-based object.
func AVOptSet(obj unsafe.Pointer, name string, val string, searchFlags int32) int32 {
	nameC, nameBuf := cString(name)
	defer runtime.KeepAlive(nameBuf)
	valC, valBuf := cString(val)
	defer runtime.KeepAlive(valBuf)
	return raw.AVOptSet(obj, nameC, valC, searchFlags)
}

// AVOptSetInt sets an integer option on an AVClass-based object.
func AVOptSetInt(obj unsafe.Pointer, name string, val int64, searchFlags int32) int32 {
	nameC, nameBuf := cString(name)
	defer runtime.KeepAlive(nameBuf)
	return raw.AVOptSetInt(obj, nameC, val, searchFlags)
}
