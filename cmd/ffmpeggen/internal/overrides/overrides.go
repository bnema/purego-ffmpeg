// cmd/ffmpeggen/internal/overrides/overrides.go
package overrides

// Domain groups C functions into an outbound port interface and public wrapper type.
type Domain struct {
	Name          string     // "format", "codec", etc.
	Library       string     // "libavformat", "libavcodec", etc.
	PortInterface string     // Go interface name: "FormatCAPI"
	PublicType    string     // Go public type: "FormatContext" (empty = free functions only)
	Functions     []FuncMap  // C→Go name mappings
	Accessors     []Accessor // offset-based struct field accessors
	Enums         []string   // enum names to include in this domain
}

type FuncMap struct {
	C  string // C function name: "avformat_open_input"
	Go string // Go method name: "OpenInput"
}

type Accessor struct {
	Struct   string // C struct name: "AVCodecContext"
	Field    string // C field name: "width"
	GoName   string // Go accessor name: "Width"
	Type     string // Go type: "int32"
	Offset   int    // byte offset in struct
	ReadOnly bool   // if true, no setter is generated (field is managed by FFmpeg internally)
}

// Override controls per-function generation behavior.
type Override struct {
	Name string
	Skip bool // skip generation entirely
}

var Overrides = []Override{
	{Name: "av_err2str", Skip: true},
	{Name: "av_ts2str", Skip: true},
	{Name: "av_ts2timestr", Skip: true},
	{Name: "av_fourcc2str", Skip: true},
}

// Accessor offsets were measured against FFmpeg 7.x (libavcodec 62, libavformat 62,
// libavutil 60) on x86_64 Linux. These offsets are ABI-dependent and MUST be
// regenerated when targeting different FFmpeg major versions or architectures.
// Use cmd/ffmpeggen/offsetgen/main.c to regenerate offset values.
var Domains = []Domain{
	{
		Name: "format", Library: "libavformat",
		PortInterface: "FormatCAPI", PublicType: "FormatContext",
		Functions: []FuncMap{
			{C: "avformat_alloc_context", Go: "AllocContext"},
			{C: "avformat_free_context", Go: "FreeContext"},
			{C: "avformat_open_input", Go: "OpenInput"},
			{C: "avformat_find_stream_info", Go: "FindStreamInfo"},
			{C: "av_read_frame", Go: "ReadFrame"},
			{C: "avformat_close_input", Go: "CloseInput"},
			{C: "avformat_alloc_output_context2", Go: "AllocOutputContext2"},
			{C: "avformat_new_stream", Go: "NewStream"},
			{C: "avformat_write_header", Go: "WriteHeader"},
			{C: "av_interleaved_write_frame", Go: "InterleavedWriteFrame"},
			{C: "av_write_trailer", Go: "WriteTrailer"},
		},
		Accessors: []Accessor{
			{Struct: "AVFormatContext", Field: "nb_streams", GoName: "NbStreams", Type: "uint32", Offset: 44, ReadOnly: true},       // managed by avformat_new_stream
			{Struct: "AVFormatContext", Field: "streams", GoName: "StreamsPtr", Type: "unsafe.Pointer", Offset: 48, ReadOnly: true}, // internal streams array
			{Struct: "AVFormatContext", Field: "duration", GoName: "Duration", Type: "int64", Offset: 104, ReadOnly: true},          // set by avformat_find_stream_info
			{Struct: "AVFormatContext", Field: "bit_rate", GoName: "BitRate", Type: "int64", Offset: 112, ReadOnly: true},           // set by avformat_find_stream_info
		},
	},
	{
		Name: "codec", Library: "libavcodec",
		PortInterface: "CodecCAPI", PublicType: "CodecContext",
		Functions: []FuncMap{
			{C: "avcodec_find_decoder", Go: "FindDecoder"},
			{C: "avcodec_find_encoder", Go: "FindEncoder"},
			{C: "avcodec_find_decoder_by_name", Go: "FindDecoderByName"},
			{C: "avcodec_find_encoder_by_name", Go: "FindEncoderByName"},
			{C: "avcodec_alloc_context3", Go: "AllocContext3"},
			{C: "avcodec_free_context", Go: "FreeContext"},
			{C: "avcodec_open2", Go: "Open2"},
			{C: "avcodec_send_packet", Go: "SendPacket"},
			{C: "avcodec_receive_frame", Go: "ReceiveFrame"},
			{C: "avcodec_send_frame", Go: "SendFrame"},
			{C: "avcodec_receive_packet", Go: "ReceivePacket"},
			{C: "avcodec_parameters_to_context", Go: "ParametersToContext"},
			{C: "avcodec_parameters_from_context", Go: "ParametersFromContext"},
		},
		Accessors: []Accessor{
			{Struct: "AVCodecContext", Field: "codec_type", GoName: "CodecType", Type: "int32", Offset: 12, ReadOnly: true}, // set by codec open
			{Struct: "AVCodecContext", Field: "codec_id", GoName: "CodecID", Type: "int32", Offset: 24, ReadOnly: true},     // set by codec open
			{Struct: "AVCodecContext", Field: "time_base", GoName: "TimeBase", Type: "AVRational", Offset: 84},
			{Struct: "AVCodecContext", Field: "width", GoName: "Width", Type: "int32", Offset: 112},
			{Struct: "AVCodecContext", Field: "height", GoName: "Height", Type: "int32", Offset: 116},
			{Struct: "AVCodecContext", Field: "pix_fmt", GoName: "PixelFormat", Type: "int32", Offset: 136},
			{Struct: "AVCodecContext", Field: "sample_rate", GoName: "SampleRate", Type: "int32", Offset: 344},
			{Struct: "AVCodecContext", Field: "sample_fmt", GoName: "SampleFormat", Type: "int32", Offset: 348},
			{Struct: "AVCodecContext", Field: "hw_device_ctx", GoName: "HWDeviceCtx", Type: "unsafe.Pointer", Offset: 560},
			{Struct: "AVCodecContext", Field: "hw_frames_ctx", GoName: "HWFramesCtx", Type: "unsafe.Pointer", Offset: 552},
		},
		Enums: []string{"AVCodecID"},
	},
	{
		Name: "packet", Library: "libavcodec",
		PortInterface: "PacketCAPI", PublicType: "Packet",
		Functions: []FuncMap{
			{C: "av_packet_alloc", Go: "Alloc"},
			{C: "av_packet_free", Go: "FreePtr"},
			{C: "av_packet_unref", Go: "Unref"},
			{C: "av_packet_ref", Go: "Ref"},
			{C: "av_packet_rescale_ts", Go: "RescaleTs"},
		},
	},
	{
		Name: "frame", Library: "libavutil",
		PortInterface: "FrameCAPI", PublicType: "Frame",
		Functions: []FuncMap{
			{C: "av_frame_alloc", Go: "Alloc"},
			{C: "av_frame_free", Go: "FreePtr"},
			{C: "av_frame_unref", Go: "Unref"},
			{C: "av_frame_ref", Go: "Ref"},
			{C: "av_frame_clone", Go: "Clone"},
			{C: "av_frame_get_buffer", Go: "GetBuffer"},
			{C: "av_frame_is_writable", Go: "IsWritable"},
			{C: "av_frame_make_writable", Go: "MakeWritable"},
		},
		Accessors: []Accessor{
			{Struct: "AVFrame", Field: "data", GoName: "DataPtr", Type: "unsafe.Pointer", Offset: 0, ReadOnly: true},          // managed by av_frame_get_buffer
			{Struct: "AVFrame", Field: "linesize", GoName: "LinesizePtr", Type: "unsafe.Pointer", Offset: 64, ReadOnly: true}, // managed by av_frame_get_buffer
			{Struct: "AVFrame", Field: "width", GoName: "Width", Type: "int32", Offset: 104},
			{Struct: "AVFrame", Field: "height", GoName: "Height", Type: "int32", Offset: 108},
			{Struct: "AVFrame", Field: "nb_samples", GoName: "NbSamples", Type: "int32", Offset: 112},
			{Struct: "AVFrame", Field: "format", GoName: "Format", Type: "int32", Offset: 116},
			{Struct: "AVFrame", Field: "pts", GoName: "Pts", Type: "int64", Offset: 136},
			{Struct: "AVFrame", Field: "pkt_dts", GoName: "PktDts", Type: "int64", Offset: 144},
			{Struct: "AVFrame", Field: "sample_rate", GoName: "SampleRate", Type: "int32", Offset: 180},
			{Struct: "AVFrame", Field: "hw_frames_ctx", GoName: "HWFramesCtx", Type: "unsafe.Pointer", Offset: 328},
		},
	},
	{
		// TODO(sws_scale): FFmpeg's sws_scale takes `const int srcStride[]` and
		// `int dstStride[]` — C arrays that decay to pointers. The header parser
		// currently maps these as scalar int32 parameters instead of unsafe.Pointer.
		// Callers must cast their stride slice pointers manually until the parser
		// handles array-type parameters or an override mechanism is added.
		Name: "swscale", Library: "libswscale",
		PortInterface: "SwscaleCAPI", PublicType: "SwscaleContext",
		Functions: []FuncMap{
			{C: "sws_getContext", Go: "GetContext"},
			{C: "sws_scale", Go: "Scale"},
			{C: "sws_freeContext", Go: "FreeContext"},
		},
	},
	{
		Name: "swresample", Library: "libswresample",
		PortInterface: "SwresampleCAPI", PublicType: "SwresampleContext",
		Functions: []FuncMap{
			{C: "swr_alloc", Go: "Alloc"},
			{C: "swr_init", Go: "Init"},
			{C: "swr_convert", Go: "Convert"},
			{C: "swr_free", Go: "FreePtr"},
		},
	},
	{
		Name: "dict", Library: "libavutil",
		PortInterface: "DictCAPI", PublicType: "Dictionary",
		Functions: []FuncMap{
			{C: "av_dict_get", Go: "Get"},
			{C: "av_dict_set", Go: "Set"},
			{C: "av_dict_free", Go: "FreePtr"},
			{C: "av_dict_count", Go: "Count"},
		},
	},
	{
		Name: "util", Library: "libavutil",
		PortInterface: "UtilCAPI", PublicType: "", // no wrapper type, free functions only
		Functions: []FuncMap{
			{C: "av_malloc", Go: "Malloc"},
			{C: "av_free", Go: "Free"},
			{C: "av_opt_set", Go: "OptSet"},
			{C: "av_opt_set_int", Go: "OptSetInt"},
			{C: "av_log_set_level", Go: "LogSetLevel"},
			{C: "av_log_get_level", Go: "LogGetLevel"},
			{C: "av_rescale_q", Go: "RescaleQ"},
			{C: "av_image_get_buffer_size", Go: "ImageGetBufferSize"},
			{C: "av_strerror", Go: "Strerror"},
		},
	},
	{
		Name: "hwaccel", Library: "libavutil",
		PortInterface: "HWAccelCAPI", PublicType: "",
		Functions: []FuncMap{
			{C: "av_hwdevice_ctx_create", Go: "DeviceCtxCreate"},
			{C: "av_hwdevice_find_type_by_name", Go: "FindTypeByName"},
			{C: "av_hwframe_transfer_data", Go: "FrameTransferData"},
			{C: "av_hwdevice_iterate_types", Go: "IterateTypes"},
		},
	},
	{
		Name: "stream", Library: "libavformat",
		PortInterface: "StreamCAPI", PublicType: "Stream",
		Functions: []FuncMap{}, // no functions — accessors only
		Accessors: []Accessor{
			{Struct: "AVStream", Field: "index", GoName: "Index", Type: "int32", Offset: 8, ReadOnly: true},                        // set by avformat
			{Struct: "AVStream", Field: "codecpar", GoName: "CodecParameters", Type: "unsafe.Pointer", Offset: 16, ReadOnly: true}, // set by avformat
			{Struct: "AVStream", Field: "time_base", GoName: "TimeBase", Type: "AVRational", Offset: 32},
			{Struct: "AVStream", Field: "duration", GoName: "Duration", Type: "int64", Offset: 48, ReadOnly: true},  // set by avformat
			{Struct: "AVStream", Field: "nb_frames", GoName: "NbFrames", Type: "int64", Offset: 56, ReadOnly: true}, // set by avformat
		},
	},
	{
		Name: "avfilter", Library: "libavfilter",
		PortInterface: "FilterCAPI", PublicType: "FilterGraph",
		Functions: []FuncMap{
			{C: "avfilter_graph_alloc", Go: "GraphAlloc"},
			{C: "avfilter_graph_free", Go: "GraphFree"},
			{C: "avfilter_graph_create_filter", Go: "GraphCreateFilter"},
			{C: "avfilter_graph_parse_ptr", Go: "GraphParsePtr"},
			{C: "avfilter_graph_config", Go: "GraphConfig"},
			{C: "avfilter_get_by_name", Go: "GetByName"},
			{C: "av_buffersrc_add_frame_flags", Go: "BuffersrcAddFrameFlags"},
			{C: "av_buffersink_get_frame", Go: "BuffersinkGetFrame"},
		},
	},
}

// TypeMap defines value types and enums shared across domains.
var Enums = []EnumDef{
	{C: "AVPixelFormat", Go: "PixelFormat"},
	{C: "AVSampleFormat", Go: "SampleFormat"},
	{C: "AVMediaType", Go: "MediaType"},
	{C: "AVCodecID", Go: "CodecID"},
	{C: "AVHWDeviceType", Go: "HWDeviceType"},
}

type EnumDef struct {
	C  string
	Go string
}

var Structs = []StructDef{
	// AVRational is NOT listed here because it appears in function signatures
	// across multiple packages (capi, out, in, ffmpeg). Each package defines it
	// via a type alias chain rooted at out.AVRational (see types.go files).
	{C: "AVDictionaryEntry", Go: "DictionaryEntry", Fields: []StructField{
		{Name: "Key", Type: "*byte"}, {Name: "Value", Type: "*byte"},
	}},
}

type StructDef struct {
	C      string
	Go     string
	Fields []StructField
}

type StructField struct {
	Name string
	Type string
}
