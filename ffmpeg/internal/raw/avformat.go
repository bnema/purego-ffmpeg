package raw

import (
	"structs"
	"unsafe"

	"github.com/ebitengine/purego"
)

type AVOutputFormat struct {
	_             structs.HostLayout
	Name          *byte
	LongName      *byte
	MimeType      *byte
	Extensions    *byte
	AudioCodec    int32
	VideoCodec    int32
	SubtitleCodec int32
	Flags         int32
	CodecTag      unsafe.Pointer
	PrivClass     unsafe.Pointer
}

type AVInputFormat struct {
	_          structs.HostLayout
	Name       *byte
	LongName   *byte
	Flags      int32
	Extensions *byte
	CodecTag   unsafe.Pointer
	PrivClass  unsafe.Pointer
	MimeType   *byte
}

// AVStream is intentionally opaque. Field access must go through the
// offset-based helpers below because the public AVStream layout changed in
// libavformat 62 and generated struct fields were no longer reliable.
type AVStream struct {
	_ structs.HostLayout
}

// ---------------------------------------------------------------------------
// AVFormatContext field offsets (FFmpeg 7.x / libavformat 62).
// ---------------------------------------------------------------------------

const (
	OffsetFmtCtxIformat   = 8
	OffsetFmtCtxOformat   = 16
	OffsetFmtCtxPB        = 32
	OffsetFmtCtxNbStreams = 44
	OffsetFmtCtxStreams   = 48
	OffsetFmtCtxURL       = 88
	OffsetFmtCtxFlags     = 128
)

// ---------------------------------------------------------------------------
// AVStream field offsets (FFmpeg 7.x / libavformat 62).
// Verified against the installed headers via offsetof().
// ---------------------------------------------------------------------------

const (
	OffsetStreamIndex             = 8
	OffsetStreamID                = 12
	OffsetStreamCodecpar          = 16
	OffsetStreamPrivData          = 24
	OffsetStreamTimeBase          = 32
	OffsetStreamStartTime         = 40
	OffsetStreamDuration          = 48
	OffsetStreamNbFrames          = 56
	OffsetStreamDisposition       = 64
	OffsetStreamDiscard           = 68
	OffsetStreamSampleAspectRatio = 72
	OffsetStreamMetadata          = 80
	OffsetStreamAvgFrameRate      = 88
	OffsetStreamAttachedPic       = 96
	OffsetStreamEventFlags        = 200
	OffsetStreamRFrameRate        = 204
	OffsetStreamPTSWrapBits       = 212
)

// ---------------------------------------------------------------------------
// Offset-based accessors for AVFormatContext
// ---------------------------------------------------------------------------

func FmtCtxPB(ctx unsafe.Pointer) unsafe.Pointer {
	if ctx == nil {
		return nil
	}
	return *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxPB))
}

func FmtCtxSetPB(ctx unsafe.Pointer, v unsafe.Pointer) {
	if ctx == nil {
		return
	}
	*(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxPB)) = v
}

func FmtCtxNbStreams(ctx unsafe.Pointer) uint32 {
	return *(*uint32)(unsafe.Add(ctx, OffsetFmtCtxNbStreams))
}

// FmtCtxStreams returns the AVStream** pointer.
func FmtCtxStreams(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxStreams))
}

// FmtCtxStream returns the i-th AVStream* from the streams array.
// Returns nil if ctx is nil or i is out of range.
func FmtCtxStream(ctx unsafe.Pointer, i int) unsafe.Pointer {
	if ctx == nil || i < 0 {
		return nil
	}
	// Bounds check against nb_streams.
	n := *(*int32)(unsafe.Add(ctx, OffsetFmtCtxNbStreams))
	if int32(i) >= n {
		return nil
	}
	arr := *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxStreams))
	// arr is AVStream**, so each element is a pointer (8 bytes on 64-bit).
	return *(*unsafe.Pointer)(unsafe.Add(arr, uintptr(i)*unsafe.Sizeof(uintptr(0))))
}

func FmtCtxOformat(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxOformat))
}

func FmtCtxSetOformat(ctx unsafe.Pointer, v unsafe.Pointer) {
	*(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxOformat)) = v
}

func FmtCtxIformat(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetFmtCtxIformat))
}

func FmtCtxFlags(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetFmtCtxFlags))
}

func FmtCtxSetFlags(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetFmtCtxFlags)) = v
}

// ---------------------------------------------------------------------------
// Offset-based accessors for AVStream
// ---------------------------------------------------------------------------

func StreamIndex(stream unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(stream, OffsetStreamIndex))
}

func StreamID(stream unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(stream, OffsetStreamID))
}

func StreamCodecpar(stream unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(stream, OffsetStreamCodecpar))
}

func StreamPrivData(stream unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(stream, OffsetStreamPrivData))
}

func StreamTimeBase(stream unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(stream, OffsetStreamTimeBase))
}

func StreamStartTime(stream unsafe.Pointer) int64 {
	return *(*int64)(unsafe.Add(stream, OffsetStreamStartTime))
}

func StreamDuration(stream unsafe.Pointer) int64 {
	return *(*int64)(unsafe.Add(stream, OffsetStreamDuration))
}

func StreamNbFrames(stream unsafe.Pointer) int64 {
	return *(*int64)(unsafe.Add(stream, OffsetStreamNbFrames))
}

func StreamDisposition(stream unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(stream, OffsetStreamDisposition))
}

func StreamDiscard(stream unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(stream, OffsetStreamDiscard))
}

func StreamSampleAspectRatio(stream unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(stream, OffsetStreamSampleAspectRatio))
}

func StreamMetadata(stream unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(stream, OffsetStreamMetadata))
}

func StreamAvgFrameRate(stream unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(stream, OffsetStreamAvgFrameRate))
}

func StreamAttachedPic(stream unsafe.Pointer) [104]byte {
	return *(*[104]byte)(unsafe.Add(stream, OffsetStreamAttachedPic))
}

func StreamEventFlags(stream unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(stream, OffsetStreamEventFlags))
}

func StreamRFrameRate(stream unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(stream, OffsetStreamRFrameRate))
}

func StreamPTSWrapBits(stream unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(stream, OffsetStreamPTSWrapBits))
}

// ---------------------------------------------------------------------------
// Function symbols
// ---------------------------------------------------------------------------

var AvformatVersion func() uint32

var AvformatAllocContext func() unsafe.Pointer

var AvformatFreeContext func(unsafe.Pointer)

var AvformatNewStream func(unsafe.Pointer, unsafe.Pointer) unsafe.Pointer

var AvformatAllocOutputContext2 func(unsafe.Pointer, unsafe.Pointer, *byte, *byte) int32

var AvformatOpenInput func(unsafe.Pointer, *byte, unsafe.Pointer, unsafe.Pointer) int32

var AvformatFindStreamInfo func(unsafe.Pointer, unsafe.Pointer) int32

var AVFindBestStream func(unsafe.Pointer, int32, int32, int32, unsafe.Pointer, int32) int32

var AVReadFrame func(unsafe.Pointer, unsafe.Pointer) int32

var AvformatCloseInput func(unsafe.Pointer)

var AvformatWriteHeader func(unsafe.Pointer, unsafe.Pointer) int32

var AVInterleavedWriteFrame func(unsafe.Pointer, unsafe.Pointer) int32

var AVWriteTrailer func(unsafe.Pointer) int32

func RegisterAvformat(handle uintptr) {
	purego.RegisterLibFunc(&AvformatVersion, handle, "avformat_version")
	purego.RegisterLibFunc(&AvformatAllocContext, handle, "avformat_alloc_context")
	purego.RegisterLibFunc(&AvformatFreeContext, handle, "avformat_free_context")
	purego.RegisterLibFunc(&AvformatNewStream, handle, "avformat_new_stream")
	purego.RegisterLibFunc(&AvformatAllocOutputContext2, handle, "avformat_alloc_output_context2")
	purego.RegisterLibFunc(&AvformatOpenInput, handle, "avformat_open_input")
	purego.RegisterLibFunc(&AvformatFindStreamInfo, handle, "avformat_find_stream_info")
	purego.RegisterLibFunc(&AVFindBestStream, handle, "av_find_best_stream")
	purego.RegisterLibFunc(&AVReadFrame, handle, "av_read_frame")
	purego.RegisterLibFunc(&AvformatCloseInput, handle, "avformat_close_input")
	purego.RegisterLibFunc(&AvformatWriteHeader, handle, "avformat_write_header")
	purego.RegisterLibFunc(&AVInterleavedWriteFrame, handle, "av_interleaved_write_frame")
	purego.RegisterLibFunc(&AVWriteTrailer, handle, "av_write_trailer")
}
