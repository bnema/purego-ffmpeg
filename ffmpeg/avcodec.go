package ffmpeg

import (
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// ---------------------------------------------------------------------------
// Codec ID constants (enum AVCodecID)
// ---------------------------------------------------------------------------

const (
	AV_CODEC_ID_H264 int32 = 27
	AV_CODEC_ID_HEVC int32 = 173
	AV_CODEC_ID_VP8  int32 = 139
	AV_CODEC_ID_VP9  int32 = 167
	AV_CODEC_ID_AV1  int32 = 225 // libavcodec/codec_id.h

	AV_CODEC_ID_AAC    int32 = 86018
	AV_CODEC_ID_OPUS   int32 = 86076
	AV_CODEC_ID_VORBIS int32 = 86021
)

// ---------------------------------------------------------------------------
// AVCodecContext accessors (offset-based, compatible with FFmpeg 7.x / libavcodec 62)
// ---------------------------------------------------------------------------

// CodecCtxCodecType returns the codec_type field of an AVCodecContext.
func CodecCtxCodecType(ctx unsafe.Pointer) int32 { return raw.CodecCtxCodecType(ctx) }

// CodecCtxSetCodecType sets the codec_type field.
func CodecCtxSetCodecType(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetCodecType(ctx, v) }

// CodecCtxCodecID returns the codec_id field of an AVCodecContext.
func CodecCtxCodecID(ctx unsafe.Pointer) int32 { return raw.CodecCtxCodecID(ctx) }

// CodecCtxSetCodecID sets the codec_id field.
func CodecCtxSetCodecID(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetCodecID(ctx, v) }

// CodecCtxWidth returns the width field of an AVCodecContext.
func CodecCtxWidth(ctx unsafe.Pointer) int32 { return raw.CodecCtxWidth(ctx) }

// CodecCtxSetWidth sets the width field.
func CodecCtxSetWidth(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetWidth(ctx, v) }

// CodecCtxHeight returns the height field of an AVCodecContext.
func CodecCtxHeight(ctx unsafe.Pointer) int32 { return raw.CodecCtxHeight(ctx) }

// CodecCtxSetHeight sets the height field.
func CodecCtxSetHeight(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetHeight(ctx, v) }

// CodecCtxPixFmt returns the pix_fmt field of an AVCodecContext.
func CodecCtxPixFmt(ctx unsafe.Pointer) int32 { return raw.CodecCtxPixFmt(ctx) }

// CodecCtxSetPixFmt sets the pix_fmt field.
func CodecCtxSetPixFmt(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetPixFmt(ctx, v) }

// CodecCtxSampleFmt returns the sample_fmt field of an AVCodecContext.
func CodecCtxSampleFmt(ctx unsafe.Pointer) int32 { return raw.CodecCtxSampleFmt(ctx) }

// CodecCtxSetSampleFmt sets the sample_fmt field.
func CodecCtxSetSampleFmt(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetSampleFmt(ctx, v) }

// CodecCtxSampleRate returns the sample_rate field of an AVCodecContext.
func CodecCtxSampleRate(ctx unsafe.Pointer) int32 { return raw.CodecCtxSampleRate(ctx) }

// CodecCtxSetSampleRate sets the sample_rate field.
func CodecCtxSetSampleRate(ctx unsafe.Pointer, v int32) { raw.CodecCtxSetSampleRate(ctx, v) }

// CodecCtxTimeBase returns the time_base field of an AVCodecContext.
func CodecCtxTimeBase(ctx unsafe.Pointer) AVRational { return raw.CodecCtxTimeBase(ctx) }

// CodecCtxSetTimeBase sets the time_base field.
func CodecCtxSetTimeBase(ctx unsafe.Pointer, v AVRational) { raw.CodecCtxSetTimeBase(ctx, v) }

// CodecCtxFramerate returns the framerate field of an AVCodecContext.
func CodecCtxFramerate(ctx unsafe.Pointer) AVRational { return raw.CodecCtxFramerate(ctx) }

// CodecCtxSetFramerate sets the framerate field.
func CodecCtxSetFramerate(ctx unsafe.Pointer, v AVRational) { raw.CodecCtxSetFramerate(ctx, v) }

// CodecCtxHwDeviceCtx returns the hw_device_ctx (AVBufferRef*) field.
func CodecCtxHwDeviceCtx(ctx unsafe.Pointer) unsafe.Pointer { return raw.CodecCtxHwDeviceCtx(ctx) }

// CodecCtxSetHwDeviceCtx sets the hw_device_ctx field.
func CodecCtxSetHwDeviceCtx(ctx unsafe.Pointer, v unsafe.Pointer) {
	raw.CodecCtxSetHwDeviceCtx(ctx, v)
}

// CodecCtxHwFramesCtx returns the hw_frames_ctx (AVBufferRef*) field.
func CodecCtxHwFramesCtx(ctx unsafe.Pointer) unsafe.Pointer { return raw.CodecCtxHwFramesCtx(ctx) }

// CodecCtxSetHwFramesCtx sets the hw_frames_ctx field.
func CodecCtxSetHwFramesCtx(ctx unsafe.Pointer, v unsafe.Pointer) {
	raw.CodecCtxSetHwFramesCtx(ctx, v)
}

// ---------------------------------------------------------------------------
// AVCodecParameters accessors
// ---------------------------------------------------------------------------

// CodecParCodecType returns the codec_type field of an AVCodecParameters.
func CodecParCodecType(par unsafe.Pointer) int32 { return raw.CodecParCodecType(par) }

// CodecParCodecID returns the codec_id field of an AVCodecParameters.
func CodecParCodecID(par unsafe.Pointer) int32 { return raw.CodecParCodecID(par) }

// CodecParWidth returns the width field of an AVCodecParameters.
func CodecParWidth(par unsafe.Pointer) int32 { return raw.CodecParWidth(par) }

// CodecParHeight returns the height field of an AVCodecParameters.
func CodecParHeight(par unsafe.Pointer) int32 { return raw.CodecParHeight(par) }

// ---------------------------------------------------------------------------
// Function wrappers
// ---------------------------------------------------------------------------

// CodecAllocContext3 allocates an AVCodecContext and sets its fields to
// default values.
func CodecAllocContext3(codec unsafe.Pointer) unsafe.Pointer {
	return raw.AvcodecAllocContext3(codec)
}

func CodecFreeContext(avctx unsafe.Pointer) {
	raw.AvcodecFreeContext(avctx)
}

// CodecParametersFromContext fills the parameters struct based on the
// values from the supplied codec context.
func CodecParametersFromContext(par unsafe.Pointer, codec unsafe.Pointer) int32 {
	return raw.AvcodecParametersFromContext(par, codec)
}

// CodecParametersToContext fills the codec context based on the values
// from the supplied codec parameters.
func CodecParametersToContext(codec unsafe.Pointer, par unsafe.Pointer) int32 {
	return raw.AvcodecParametersToContext(codec, par)
}

// CodecParametersCopy copies codec parameters from src to dst.
func CodecParametersCopy(dst unsafe.Pointer, src unsafe.Pointer) int32 {
	return raw.AvcodecParametersCopy(dst, src)
}

func CodecOpen2(avctx unsafe.Pointer, codec unsafe.Pointer, options unsafe.Pointer) int32 {
	return raw.AvcodecOpen2(avctx, codec, options)
}

func CodecSendPacket(avctx unsafe.Pointer, avpkt unsafe.Pointer) int32 {
	return raw.AvcodecSendPacket(avctx, avpkt)
}

func CodecReceiveFrame(avctx unsafe.Pointer, frame unsafe.Pointer) int32 {
	return raw.AvcodecReceiveFrame(avctx, frame)
}

func CodecSendFrame(avctx unsafe.Pointer, frame unsafe.Pointer) int32 {
	return raw.AvcodecSendFrame(avctx, frame)
}

func CodecReceivePacket(avctx unsafe.Pointer, avpkt unsafe.Pointer) int32 {
	return raw.AvcodecReceivePacket(avctx, avpkt)
}

func CodecFlushBuffers(avctx unsafe.Pointer) {
	raw.AvcodecFlushBuffers(avctx)
}

// CodecIsOpen returns non-zero if the codec has been opened.
func CodecIsOpen(s unsafe.Pointer) int32 {
	return raw.AvcodecIsOpen(s)
}
