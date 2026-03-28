package ffmpeg

import (
	"unsafe"

	"github.com/bnema/purego-ffmpeg/internal/capi"
	"github.com/bnema/purego-ffmpeg/internal/core"
)

// ------------------------------------------------------------
// Codec package-level functions
// ------------------------------------------------------------

func CodecFindDecoder(id int32) unsafe.Pointer {
	return defaultCodec().FindDecoder(id)
}

func CodecFindEncoder(id int32) unsafe.Pointer {
	return defaultCodec().FindEncoder(id)
}

func CodecFindDecoderByName(name string) unsafe.Pointer {
	cName := core.CString(name)
	defer core.FreeCString(cName)
	return defaultCodec().FindDecoderByName(cName)
}

func CodecFindEncoderByName(name string) unsafe.Pointer {
	cName := core.CString(name)
	defer core.FreeCString(cName)
	return defaultCodec().FindEncoderByName(cName)
}

func CodecAllocContext3(codec unsafe.Pointer) unsafe.Pointer {
	return defaultCodec().AllocContext3(codec)
}

func CodecFreeContext(avctx unsafe.Pointer) {
	defaultCodec().FreeContext(avctx)
}

func CodecOpen2(avctx, codec, options unsafe.Pointer) int32 {
	return defaultCodec().Open2(avctx, codec, options)
}

func CodecParametersToContext(codec, par unsafe.Pointer) int32 {
	return defaultCodec().ParametersToContext(codec, par)
}

func CodecParametersFromContext(par, codec unsafe.Pointer) int32 {
	return defaultCodec().ParametersFromContext(par, codec)
}

func CodecSendPacket(avctx, avpkt unsafe.Pointer) int32 {
	return defaultCodec().SendPacket(avctx, avpkt)
}

func CodecReceiveFrame(avctx, frame unsafe.Pointer) int32 {
	return defaultCodec().ReceiveFrame(avctx, frame)
}

func CodecSendFrame(avctx, frame unsafe.Pointer) int32 {
	return defaultCodec().SendFrame(avctx, frame)
}

func CodecReceivePacket(avctx, avpkt unsafe.Pointer) int32 {
	return defaultCodec().ReceivePacket(avctx, avpkt)
}

// ------------------------------------------------------------
// Format package-level functions
// ------------------------------------------------------------

func FormatAllocContext() unsafe.Pointer {
	return defaultFormat().AllocContext()
}

func FormatFreeContext(s unsafe.Pointer) {
	defaultFormat().FreeContext(s)
}

func FormatOpenInput(ps unsafe.Pointer, url string, fmt, options unsafe.Pointer) int32 {
	cURL := core.CString(url)
	defer core.FreeCString(cURL)
	return defaultFormat().OpenInput(ps, cURL, fmt, options)
}

func FormatCloseInput(s unsafe.Pointer) {
	defaultFormat().CloseInput(s)
}

func FormatFindStreamInfo(ic, options unsafe.Pointer) int32 {
	return defaultFormat().FindStreamInfo(ic, options)
}

func FormatAllocOutputContext2(ctx, oformat unsafe.Pointer, formatName, filename string) int32 {
	cFmt := core.CString(formatName)
	defer core.FreeCString(cFmt)
	cFile := core.CString(filename)
	defer core.FreeCString(cFile)
	return defaultFormat().AllocOutputContext2(ctx, oformat, cFmt, cFile)
}

func FormatNewStream(s, c unsafe.Pointer) unsafe.Pointer {
	return defaultFormat().NewStream(s, c)
}

func FormatWriteHeader(s, options unsafe.Pointer) int32 {
	return defaultFormat().WriteHeader(s, options)
}

func InterleavedWriteFrame(s, pkt unsafe.Pointer) int32 {
	return defaultFormat().InterleavedWriteFrame(s, pkt)
}

func WriteTrailer(s unsafe.Pointer) int32 {
	return defaultFormat().WriteTrailer(s)
}

func ReadFrame(s, pkt unsafe.Pointer) int32 {
	return defaultFormat().ReadFrame(s, pkt)
}

// ------------------------------------------------------------
// Frame / Packet raw-pointer functions
// ------------------------------------------------------------

func FrameAlloc() unsafe.Pointer {
	return defaultFrame().Alloc()
}

func FrameFree(frame unsafe.Pointer) {
	defaultFrame().FreePtr(frame)
}

func FrameUnref(frame unsafe.Pointer) {
	defaultFrame().Unref(frame)
}

func FrameGetBuffer(frame unsafe.Pointer, align int32) int32 {
	return defaultFrame().GetBuffer(frame, align)
}

func PacketAlloc() unsafe.Pointer {
	return defaultPacket().Alloc()
}

func PacketFree(pkt unsafe.Pointer) {
	defaultPacket().FreePtr(pkt)
}

func PacketUnref(pkt unsafe.Pointer) {
	defaultPacket().Unref(pkt)
}

func PacketRescaleTs(pkt unsafe.Pointer, tbSrc, tbDst AVRational) {
	defaultPacket().RescaleTs(pkt, tbSrc, tbDst)
}

// ------------------------------------------------------------
// CodecContext field accessors
// ------------------------------------------------------------

func CodecCtxWidth(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextWidth))
}

func CodecCtxSetWidth(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextWidth)) = v
}

func CodecCtxHeight(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextHeight))
}

func CodecCtxSetHeight(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextHeight)) = v
}

func CodecCtxPixFmt(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextPixelFormat))
}

func CodecCtxSetPixFmt(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextPixelFormat)) = v
}

func CodecCtxSampleRate(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextSampleRate))
}

func CodecCtxSetSampleRate(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextSampleRate)) = v
}

func CodecCtxSampleFmt(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextSampleFormat))
}

func CodecCtxSetSampleFmt(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextSampleFormat)) = v
}

func CodecCtxTimeBase(ctx unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(ctx, capi.OffsetAVCodecContextTimeBase))
}

func CodecCtxSetTimeBase(ctx unsafe.Pointer, v AVRational) {
	*(*AVRational)(unsafe.Add(ctx, capi.OffsetAVCodecContextTimeBase)) = v
}

func CodecCtxFramerate(ctx unsafe.Pointer) AVRational {
	return *(*AVRational)(unsafe.Add(ctx, capi.OffsetAVCodecContextFramerate))
}

func CodecCtxSetFramerate(ctx unsafe.Pointer, v AVRational) {
	*(*AVRational)(unsafe.Add(ctx, capi.OffsetAVCodecContextFramerate)) = v
}

func CodecCtxBitRate(ctx unsafe.Pointer) int64 {
	return *(*int64)(unsafe.Add(ctx, capi.OffsetAVCodecContextBitRate))
}

func CodecCtxSetBitRate(ctx unsafe.Pointer, v int64) {
	*(*int64)(unsafe.Add(ctx, capi.OffsetAVCodecContextBitRate)) = v
}

func CodecCtxGopSize(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextGopSize))
}

func CodecCtxSetGopSize(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextGopSize)) = v
}

func CodecCtxFlags(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextFlags))
}

func CodecCtxSetFlags(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVCodecContextFlags)) = v
}

func CodecCtxHwDeviceCtx(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, capi.OffsetAVCodecContextHWDeviceCtx))
}

func CodecCtxSetHwDeviceCtx(ctx unsafe.Pointer, v unsafe.Pointer) {
	*(*unsafe.Pointer)(unsafe.Add(ctx, capi.OffsetAVCodecContextHWDeviceCtx)) = v
}

func CodecCtxHwFramesCtx(ctx unsafe.Pointer) unsafe.Pointer {
	return *(*unsafe.Pointer)(unsafe.Add(ctx, capi.OffsetAVCodecContextHWFramesCtx))
}

func CodecCtxSetHwFramesCtx(ctx unsafe.Pointer, v unsafe.Pointer) {
	*(*unsafe.Pointer)(unsafe.Add(ctx, capi.OffsetAVCodecContextHWFramesCtx)) = v
}

// ------------------------------------------------------------
// CodecParameters accessors
// ------------------------------------------------------------

func CodecParCodecID(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, capi.OffsetAVCodecParametersCodecID))
}

func CodecParCodecType(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, capi.OffsetAVCodecParametersCodecType))
}

func CodecParWidth(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, capi.OffsetAVCodecParametersWidth))
}

func CodecParHeight(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, capi.OffsetAVCodecParametersHeight))
}

func CodecParSampleRate(par unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(par, capi.OffsetAVCodecParametersSampleRate))
}

// ------------------------------------------------------------
// Format context field accessors
// ------------------------------------------------------------

func FmtCtxNbStreams(ctx unsafe.Pointer) uint32 {
	return *(*uint32)(unsafe.Add(ctx, capi.OffsetAVFormatContextNbStreams))
}

// FmtCtxStream returns the AVStream* at index idx from AVFormatContext.streams[idx].
// Returns nil if idx is out of range.
func FmtCtxStream(ctx unsafe.Pointer, idx int) unsafe.Pointer {
	if idx < 0 || uint32(idx) >= FmtCtxNbStreams(ctx) {
		return nil
	}
	streamsPtr := *(*unsafe.Pointer)(unsafe.Add(ctx, capi.OffsetAVFormatContextStreamsPtr))
	if streamsPtr == nil {
		return nil
	}
	// streams is AVStream** — array of pointers, each 8 bytes on 64-bit
	return *(*unsafe.Pointer)(unsafe.Add(streamsPtr, uintptr(idx)*unsafe.Sizeof(uintptr(0))))
}

func FmtCtxSetPB(ctx unsafe.Pointer, pb unsafe.Pointer) {
	*(*unsafe.Pointer)(unsafe.Add(ctx, capi.OffsetAVFormatContextPB)) = pb
}

func FmtCtxFlags(ctx unsafe.Pointer) int32 {
	return *(*int32)(unsafe.Add(ctx, capi.OffsetAVFormatContextFlags))
}

func FmtCtxSetFlags(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVFormatContextFlags)) = v
}

// ------------------------------------------------------------
// AVIO functions
// ------------------------------------------------------------

// AvioAllocContext allocates an AVIOContext for custom I/O.
// readCb/writeCb/seekCb are purego function pointers (use purego.NewCallback).
// The internal buffer is allocated by FFmpeg (nil is passed for the buffer parameter).
func AvioAllocContext(bufSize int, writeFlag bool, opaque unsafe.Pointer, readCb, writeCb, seekCb uintptr) unsafe.Pointer {
	wf := int32(0)
	if writeFlag {
		wf = 1
	}
	return defaultFormat().AvioAllocContext(nil, int32(bufSize), wf, opaque, readCb, writeCb, seekCb)
}

func AvioContextFree(ctx unsafe.Pointer) {
	defaultFormat().AvioContextFree(ctx)
}

// ------------------------------------------------------------
// Buffer reference counting
// ------------------------------------------------------------

func BufferRef(buf unsafe.Pointer) unsafe.Pointer {
	return defaultUtil().BufferRef(buf)
}

func BufferUnref(buf *unsafe.Pointer) {
	if buf != nil && *buf != nil {
		defaultUtil().BufferUnref(unsafe.Pointer(buf))
		*buf = nil
	}
}

// BufferRefData returns the data pointer from an AVBufferRef.
func BufferRefData(buf unsafe.Pointer) unsafe.Pointer {
	if buf == nil {
		return nil
	}
	return *(*unsafe.Pointer)(unsafe.Add(buf, capi.OffsetAVBufferRefData))
}

// ------------------------------------------------------------
// HW Frames Context
// ------------------------------------------------------------

func HWFrameCtxAlloc(deviceCtx unsafe.Pointer) unsafe.Pointer {
	return defaultHWAccel().HWFrameCtxAlloc(deviceCtx)
}

func HWFrameCtxInit(ref unsafe.Pointer) int32 {
	return defaultHWAccel().HWFrameCtxInit(ref)
}

func HWFrameGetBuffer(hwFramesCtx, frame unsafe.Pointer, flags int32) int32 {
	return defaultHWAccel().HWFrameGetBuffer(hwFramesCtx, frame, flags)
}

// HWFrameTransferData already exists as FrameTransferData in hwaccel_gen.go.
// Re-export for consumer compatibility.
func HWFrameTransferData(dst, src unsafe.Pointer, flags int32) int32 {
	return FrameTransferData(dst, src, flags)
}

// HWFramesCtxSetFormat sets the pixel format on an AVHWFramesContext.
// ctx must point to AVHWFramesContext (obtained via BufferRefData), not AVBufferRef.
func HWFramesCtxSetFormat(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVHWFramesContextFormat)) = v
}

func HWFramesCtxSetSWFormat(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVHWFramesContextSWFormat)) = v
}

func HWFramesCtxSetWidth(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVHWFramesContextWidth)) = v
}

func HWFramesCtxSetHeight(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVHWFramesContextHeight)) = v
}

func HWFramesCtxSetInitialPoolSize(ctx unsafe.Pointer, v int32) {
	*(*int32)(unsafe.Add(ctx, capi.OffsetAVHWFramesContextInitialPoolSize)) = v
}

// ------------------------------------------------------------
// Convenience / Utility
// ------------------------------------------------------------

// HWDeviceCtxCreate is a Go-friendly wrapper around DeviceCtxCreate.
func HWDeviceCtxCreate(deviceType int32, device string) (unsafe.Pointer, error) {
	var deviceCtx unsafe.Pointer
	cDevice := core.CString(device)
	defer core.FreeCString(cDevice)
	ret := DeviceCtxCreate(unsafe.Pointer(&deviceCtx), deviceType, cDevice, nil, 0)
	if err := core.CheckError(ret); err != nil {
		return nil, err
	}
	return deviceCtx, nil
}

// DictSet wraps av_dict_set with Go strings.
func DictSet(dict unsafe.Pointer, key, value string, flags int32) int32 {
	cKey := core.CString(key)
	defer core.FreeCString(cKey)
	cVal := core.CString(value)
	defer core.FreeCString(cVal)
	return defaultDict().Set(dict, cKey, cVal, flags)
}

func DictFree(dict unsafe.Pointer) {
	defaultDict().FreePtr(dict)
}

func GetSampleFmtName(sampleFmt int32) string {
	p := defaultUtil().GetSampleFmtName(sampleFmt)
	return core.GoString(p)
}

// AVOptSet wraps av_opt_set with Go strings.
func AVOptSet(obj unsafe.Pointer, name, val string, searchFlags int32) int32 {
	cName := core.CString(name)
	defer core.FreeCString(cName)
	cVal := core.CString(val)
	defer core.FreeCString(cVal)
	return defaultUtil().OptSet(obj, cName, cVal, searchFlags)
}

// AVOptSetInt wraps av_opt_set_int with a Go string name.
func AVOptSetInt(obj unsafe.Pointer, name string, val int64, searchFlags int32) int32 {
	cName := core.CString(name)
	defer core.FreeCString(cName)
	return defaultUtil().OptSetInt(obj, cName, val, searchFlags)
}

func SwrAlloc() unsafe.Pointer {
	return defaultSwresample().Alloc()
}

func SwrInit(swr unsafe.Pointer) int32 {
	return defaultSwresample().Init(swr)
}

func SwrConvert(swr, out unsafe.Pointer, outCount int32, in unsafe.Pointer, inCount int32) int32 {
	return defaultSwresample().Convert(swr, out, outCount, in, inCount)
}

func SwrFree(swr unsafe.Pointer) {
	defaultSwresample().FreePtr(swr)
}

func SwrGetDelay(swr unsafe.Pointer, base int64) int64 {
	return defaultSwresample().GetDelay(swr, base)
}

// ------------------------------------------------------------
// Constants
// ------------------------------------------------------------

const (
	AVFMT_FLAG_CUSTOM_IO = 0x0080
	AVMEDIA_TYPE_VIDEO   = int32(AvmediaTypeVideo)
	AVMEDIA_TYPE_AUDIO   = int32(AvmediaTypeAudio)
)
