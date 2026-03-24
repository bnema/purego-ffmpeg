package raw

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// ---------------------------------------------------------------------------
// AVCodecContext field offsets (FFmpeg 7.x / libavcodec 62).
// Obtained via offsetof() against the installed headers.
// ---------------------------------------------------------------------------

const (
	OffsetCodecCtxCodecType   = 12
	OffsetCodecCtxCodecID     = 24
	OffsetCodecCtxTimeBase    = 84
	OffsetCodecCtxFramerate   = 100
	OffsetCodecCtxWidth       = 112
	OffsetCodecCtxHeight      = 116
	OffsetCodecCtxPixFmt      = 136
	OffsetCodecCtxSampleRate  = 344
	OffsetCodecCtxSampleFmt   = 348
	OffsetCodecCtxHwFramesCtx = 552
	OffsetCodecCtxHwDeviceCtx = 560
)

// ---------------------------------------------------------------------------
// AVCodecParameters field offsets (FFmpeg 7.x / libavcodec 62).
// ---------------------------------------------------------------------------

const (
	OffsetCodecParCodecType  = 0
	OffsetCodecParCodecID    = 4
	OffsetCodecParCodecTag   = 8
	OffsetCodecParFormat     = 44
	OffsetCodecParBitRate    = 48
	OffsetCodecParWidth      = 72
	OffsetCodecParHeight     = 76
	OffsetCodecParSampleRate = 152
)

// ---------------------------------------------------------------------------
// Offset-based accessors for AVCodecContext.
// All functions in this section require ctx to be non-nil. No nil guards are
// added here because these are hot-path FFI accessors; callers must ensure
// a valid context before invoking.
// ---------------------------------------------------------------------------

func CodecCtxCodecType(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxCodecType))
}

func CodecCtxSetCodecType(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxCodecType)) = v
}

func CodecCtxCodecID(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxCodecID))
}

func CodecCtxSetCodecID(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxCodecID)) = v
}

func CodecCtxWidth(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxWidth))
}

func CodecCtxSetWidth(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxWidth)) = v
}

func CodecCtxHeight(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxHeight))
}

func CodecCtxSetHeight(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxHeight)) = v
}

func CodecCtxPixFmt(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxPixFmt))
}

func CodecCtxSetPixFmt(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxPixFmt)) = v
}

func CodecCtxSampleFmt(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxSampleFmt))
}

func CodecCtxSetSampleFmt(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxSampleFmt)) = v
}

func CodecCtxSampleRate(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, OffsetCodecCtxSampleRate))
}

func CodecCtxSetSampleRate(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, OffsetCodecCtxSampleRate)) = v
}

func CodecCtxTimeBase(ctx unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(ctx, OffsetCodecCtxTimeBase))
}

func CodecCtxSetTimeBase(ctx unsafe.Pointer, v AVRational) {
	*(*AVRational)(unsafe.Add(ctx, OffsetCodecCtxTimeBase)) = v
}

func CodecCtxFramerate(ctx unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(ctx, OffsetCodecCtxFramerate))
}

func CodecCtxSetFramerate(ctx unsafe.Pointer, v AVRational) {
	*(*AVRational)(unsafe.Add(ctx, OffsetCodecCtxFramerate)) = v
}

func CodecCtxHwDeviceCtx(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetCodecCtxHwDeviceCtx))
}

func CodecCtxSetHwDeviceCtx(ctx unsafe.Pointer, v unsafe.Pointer) {
	*(*unsafe.Pointer)(unsafe.Add(ctx, OffsetCodecCtxHwDeviceCtx)) = v
}

func CodecCtxHwFramesCtx(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, OffsetCodecCtxHwFramesCtx))
}

func CodecCtxSetHwFramesCtx(ctx unsafe.Pointer, v unsafe.Pointer) {
	*(*unsafe.Pointer)(unsafe.Add(ctx, OffsetCodecCtxHwFramesCtx)) = v
}

// ---------------------------------------------------------------------------
// Offset-based accessors for AVCodecParameters
// ---------------------------------------------------------------------------

func CodecParCodecType(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, OffsetCodecParCodecType))
}

func CodecParCodecID(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, OffsetCodecParCodecID))
}

func CodecParWidth(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, OffsetCodecParWidth))
}

func CodecParHeight(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, OffsetCodecParHeight))
}

// ---------------------------------------------------------------------------
// Function symbols
// ---------------------------------------------------------------------------

var AvcodecAllocContext3 func(unsafe.Pointer) unsafe.Pointer

var AvcodecFreeContext func(unsafe.Pointer)

var AvcodecParametersFromContext func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecParametersToContext func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecParametersCopy func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecOpen2 func(unsafe.Pointer, unsafe.Pointer, unsafe.Pointer) int32

var AvcodecSendPacket func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecReceiveFrame func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecSendFrame func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecReceivePacket func(unsafe.Pointer, unsafe.Pointer) int32

var AvcodecFlushBuffers func(unsafe.Pointer)

var AvcodecIsOpen func(unsafe.Pointer) int32

func RegisterAvcodec(handle uintptr) {
	purego.RegisterLibFunc(&AvcodecAllocContext3, handle, "avcodec_alloc_context3")
	purego.RegisterLibFunc(&AvcodecFreeContext, handle, "avcodec_free_context")
	purego.RegisterLibFunc(&AvcodecParametersFromContext, handle, "avcodec_parameters_from_context")
	purego.RegisterLibFunc(&AvcodecParametersToContext, handle, "avcodec_parameters_to_context")
	purego.RegisterLibFunc(&AvcodecParametersCopy, handle, "avcodec_parameters_copy")
	purego.RegisterLibFunc(&AvcodecOpen2, handle, "avcodec_open2")
	purego.RegisterLibFunc(&AvcodecSendPacket, handle, "avcodec_send_packet")
	purego.RegisterLibFunc(&AvcodecReceiveFrame, handle, "avcodec_receive_frame")
	purego.RegisterLibFunc(&AvcodecSendFrame, handle, "avcodec_send_frame")
	purego.RegisterLibFunc(&AvcodecReceivePacket, handle, "avcodec_receive_packet")
	purego.RegisterLibFunc(&AvcodecFlushBuffers, handle, "avcodec_flush_buffers")
	purego.RegisterLibFunc(&AvcodecIsOpen, handle, "avcodec_is_open")
}
