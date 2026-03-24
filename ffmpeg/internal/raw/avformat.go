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

type AVStream struct {
	_                 structs.HostLayout
	AvClass           unsafe.Pointer
	Index             int32
	ID                int32
	Codecpar          unsafe.Pointer
	PrivData          unsafe.Pointer
	TimeBase          AVRational
	StartTime         int64
	Duration          int64
	NbFrames          int64
	Disposition       int32
	Discard           int32
	SampleAspectRatio AVRational
	Metadata          unsafe.Pointer
	AvgFrameRate      AVRational
	AttachedPic       [104]byte
	EventFlags        int32
	RFrameRate        AVRational
	PTSWrapBits       int32
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
