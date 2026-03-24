package ffmpeg

import (
	"runtime"
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// OutputFormat wraps an AVOutputFormat.
type OutputFormat struct {
	ptr *raw.AVOutputFormat
}

func (w *OutputFormat) Raw() unsafe.Pointer {
	if w == nil {
		return nil
	}
	return unsafe.Pointer(w.ptr)
}

func (w *OutputFormat) Name() *byte {
	return w.ptr.Name
}

func (w *OutputFormat) LongName() *byte {
	return w.ptr.LongName
}

func (w *OutputFormat) MimeType() *byte {
	return w.ptr.MimeType
}

func (w *OutputFormat) Extensions() *byte {
	return w.ptr.Extensions
}

func (w *OutputFormat) AudioCodec() int32 {
	return w.ptr.AudioCodec
}

func (w *OutputFormat) VideoCodec() int32 {
	return w.ptr.VideoCodec
}

func (w *OutputFormat) SubtitleCodec() int32 {
	return w.ptr.SubtitleCodec
}

func (w *OutputFormat) Flags() int32 {
	return w.ptr.Flags
}

func (w *OutputFormat) CodecTag() unsafe.Pointer {
	return w.ptr.CodecTag
}

func (w *OutputFormat) PrivClass() unsafe.Pointer {
	return w.ptr.PrivClass
}

// InputFormat wraps an AVInputFormat.
type InputFormat struct {
	ptr *raw.AVInputFormat
}

func (w *InputFormat) Raw() unsafe.Pointer {
	if w == nil {
		return nil
	}
	return unsafe.Pointer(w.ptr)
}

func (w *InputFormat) Name() *byte {
	return w.ptr.Name
}

func (w *InputFormat) LongName() *byte {
	return w.ptr.LongName
}

func (w *InputFormat) Flags() int32 {
	return w.ptr.Flags
}

func (w *InputFormat) Extensions() *byte {
	return w.ptr.Extensions
}

func (w *InputFormat) CodecTag() unsafe.Pointer {
	return w.ptr.CodecTag
}

func (w *InputFormat) PrivClass() unsafe.Pointer {
	return w.ptr.PrivClass
}

func (w *InputFormat) MimeType() *byte {
	return w.ptr.MimeType
}

// Stream wraps an AVStream.
type Stream struct {
	ptr *raw.AVStream
}

func (w *Stream) Raw() unsafe.Pointer {
	if w == nil {
		return nil
	}
	return unsafe.Pointer(w.ptr)
}

func (w *Stream) AvClass() unsafe.Pointer {
	return w.ptr.AvClass
}

func (w *Stream) Index() int32 {
	return w.ptr.Index
}

func (w *Stream) ID() int32 {
	return w.ptr.ID
}

func (w *Stream) Codecpar() unsafe.Pointer {
	return w.ptr.Codecpar
}

func (w *Stream) PrivData() unsafe.Pointer {
	return w.ptr.PrivData
}

func (w *Stream) TimeBase() AVRational {
	return w.ptr.TimeBase
}

func (w *Stream) StartTime() int64 {
	return w.ptr.StartTime
}

func (w *Stream) Duration() int64 {
	return w.ptr.Duration
}

func (w *Stream) NbFrames() int64 {
	return w.ptr.NbFrames
}

func (w *Stream) Disposition() int32 {
	return w.ptr.Disposition
}

func (w *Stream) Discard() int32 {
	return w.ptr.Discard
}

func (w *Stream) SampleAspectRatio() AVRational {
	return w.ptr.SampleAspectRatio
}

func (w *Stream) Metadata() unsafe.Pointer {
	return w.ptr.Metadata
}

func (w *Stream) AvgFrameRate() AVRational {
	return w.ptr.AvgFrameRate
}

func (w *Stream) AttachedPic() [104]byte {
	return w.ptr.AttachedPic
}

func (w *Stream) EventFlags() int32 {
	return w.ptr.EventFlags
}

func (w *Stream) RFrameRate() AVRational {
	return w.ptr.RFrameRate
}

func (w *Stream) PTSWrapBits() int32 {
	return w.ptr.PTSWrapBits
}

// ---------------------------------------------------------------------------
// AVFormatContext accessors (offset-based)
// ---------------------------------------------------------------------------

// FmtCtxPB returns the pb (AVIOContext*) field of an AVFormatContext.
func FmtCtxPB(ctx unsafe.Pointer) unsafe.Pointer { return raw.FmtCtxPB(ctx) }

// FmtCtxSetPB sets the pb field to a custom AVIO context.
func FmtCtxSetPB(ctx unsafe.Pointer, pb unsafe.Pointer) { raw.FmtCtxSetPB(ctx, pb) }

// FmtCtxNbStreams returns the nb_streams field.
func FmtCtxNbStreams(ctx unsafe.Pointer) uint32 { return raw.FmtCtxNbStreams(ctx) }

// FmtCtxStreams returns the raw streams pointer (AVStream**).
func FmtCtxStreams(ctx unsafe.Pointer) unsafe.Pointer { return raw.FmtCtxStreams(ctx) }

// FmtCtxStream returns the i-th AVStream* from the format context.
func FmtCtxStream(ctx unsafe.Pointer, i int) unsafe.Pointer { return raw.FmtCtxStream(ctx, i) }

// FmtCtxOformat returns the oformat (AVOutputFormat*) field.
func FmtCtxOformat(ctx unsafe.Pointer) unsafe.Pointer { return raw.FmtCtxOformat(ctx) }

// FmtCtxSetOformat sets the oformat field.
func FmtCtxSetOformat(ctx unsafe.Pointer, v unsafe.Pointer) { raw.FmtCtxSetOformat(ctx, v) }

// FmtCtxIformat returns the iformat (AVInputFormat*) field.
func FmtCtxIformat(ctx unsafe.Pointer) unsafe.Pointer { return raw.FmtCtxIformat(ctx) }

// FmtCtxFlags returns the flags field.
func FmtCtxFlags(ctx unsafe.Pointer) int32 { return raw.FmtCtxFlags(ctx) }

// FmtCtxSetFlags sets the flags field.
func FmtCtxSetFlags(ctx unsafe.Pointer, v int32) { raw.FmtCtxSetFlags(ctx, v) }

// AVFormatContext flags (used with FmtCtxSetFlags).
const AVFMT_FLAG_CUSTOM_IO = 0x0080

// AVOutputFormat/AVInputFormat flags (format-level, not context-level).
const (
	AVFMT_NOFILE       = 0x0001
	AVFMT_GLOBALHEADER = 0x0040
)

// ---------------------------------------------------------------------------
// Function wrappers
// ---------------------------------------------------------------------------

// FormatVersion returns the LIBAVFORMAT_VERSION_INT constant.
func FormatVersion() uint32 {
	return raw.AvformatVersion()
}

func FormatAllocContext() unsafe.Pointer {
	return raw.AvformatAllocContext()
}

func FormatFreeContext(s unsafe.Pointer) {
	raw.AvformatFreeContext(s)
}

func FormatNewStream(s unsafe.Pointer, c unsafe.Pointer) unsafe.Pointer {
	return raw.AvformatNewStream(s, c)
}

// FormatAllocOutputContext2 allocates an AVFormatContext for an output format.
func FormatAllocOutputContext2(ctx unsafe.Pointer, oformat unsafe.Pointer, formatName string, filename string) int32 {
	formatNameC, formatNameBuf := cString(formatName)
	defer runtime.KeepAlive(formatNameBuf)
	filenameC, filenameBuf := cString(filename)
	defer runtime.KeepAlive(filenameBuf)
	return raw.AvformatAllocOutputContext2(ctx, oformat, formatNameC, filenameC)
}

func FormatOpenInput(ps unsafe.Pointer, uRL string, fmt unsafe.Pointer, options unsafe.Pointer) int32 {
	uRLC, uRLBuf := cString(uRL)
	defer runtime.KeepAlive(uRLBuf)
	return raw.AvformatOpenInput(ps, uRLC, fmt, options)
}

func FormatFindStreamInfo(ic unsafe.Pointer, options unsafe.Pointer) int32 {
	return raw.AvformatFindStreamInfo(ic, options)
}

// FindBestStream finds the "best" stream in the file.
func FindBestStream(ic unsafe.Pointer, type_ int32, wantedStreamNb int32, relatedStream int32, decoderRet unsafe.Pointer, flags int32) int32 {
	return raw.AVFindBestStream(ic, type_, wantedStreamNb, relatedStream, decoderRet, flags)
}

func ReadFrame(s unsafe.Pointer, pkt unsafe.Pointer) int32 {
	return raw.AVReadFrame(s, pkt)
}

func FormatCloseInput(s unsafe.Pointer) {
	raw.AvformatCloseInput(s)
}

func FormatWriteHeader(s unsafe.Pointer, options unsafe.Pointer) int32 {
	return raw.AvformatWriteHeader(s, options)
}

func InterleavedWriteFrame(s unsafe.Pointer, pkt unsafe.Pointer) int32 {
	return raw.AVInterleavedWriteFrame(s, pkt)
}

func WriteTrailer(s unsafe.Pointer) int32 {
	return raw.AVWriteTrailer(s)
}
