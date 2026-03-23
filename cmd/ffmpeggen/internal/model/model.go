package model

import "strings"

// Header represents a parsed FFmpeg C header file.
type Header struct {
	Path         string
	Lib          string // "avutil", "avformat", "avcodec", "swscale", "swresample"
	RegisterName string // e.g. "RegisterAvformat"
	Structs      []Struct
	Functions    []Function
	Enums        []Enum
}

// Struct represents a parsed C struct.
type Struct struct {
	CName      string
	GoName     string
	PublicName string // clean public Go name (e.g., "FormatContext")
	Kind       string // "opaque", "data", "config"
	Doc        string
	Fields     []Field
	IsOpaque   bool // true if struct is forward-declared only
}

// NeedsUnsafe reports whether the header uses unsafe.Pointer.
func (h *Header) NeedsUnsafe() bool {
	for _, s := range h.Structs {
		for _, f := range s.Fields {
			if f.GoType == "unsafe.Pointer" {
				return true
			}
		}
	}
	for _, fn := range h.Functions {
		if fn.ReturnGoType == "unsafe.Pointer" {
			return true
		}
		for _, p := range fn.Params {
			if p.GoType == "unsafe.Pointer" {
				return true
			}
		}
	}
	return false
}

// NeedsPurego reports whether the header has functions requiring purego.
func (h *Header) NeedsPurego() bool {
	return len(h.Functions) > 0
}

// NeedsStructs reports whether the header defines any structs.
func (h *Header) NeedsStructs() bool {
	return len(h.Structs) > 0
}

// Field represents a single field in a C struct.
type Field struct {
	CName  string
	GoName string
	CType  string
	GoType string
	Doc    string
	// Array info for fixed-size C arrays (e.g., uint8_t data[AV_NUM_DATA_POINTERS])
	IsArray   bool
	ArraySize string // Go expression for array length
}

// Function represents a parsed C function declaration.
type Function struct {
	CName        string
	GoName       string
	Doc          string
	Params       []Param
	ReturnCType  string
	ReturnGoType string
}

// Param represents a single function parameter.
type Param struct {
	CName   string
	GoName  string
	CType   string
	GoType  string
	IsConst bool
	Pointer int
}

// Enum represents a parsed C enum definition.
type Enum struct {
	CName  string
	GoName string
	Doc    string
	Values []EnumValue
}

// EnumValue represents a single value in a C enum.
type EnumValue struct {
	CName  string
	GoName string
	Value  string
}

// acronyms maps lowercase segments to their Go equivalents.
var acronyms = map[string]string{
	"id":   "ID",
	"url":  "URL",
	"io":   "IO",
	"eof":  "EOF",
	"pts":  "PTS",
	"dts":  "DTS",
	"fps":  "FPS",
	"rgb":  "RGB",
	"yuv":  "YUV",
	"cpu":  "CPU",
	"gpu":  "GPU",
	"hw":   "HW",
	"api":  "API",
	"aac":  "AAC",
	"mp3":  "MP3",
	"h264": "H264",
	"h265": "H265",
	"av":   "AV",
	"sws":  "SWS",
	"swr":  "SWR",
}

// GoStructName converts a C struct name to a Go raw struct name.
// FFmpeg names are already PascalCase (e.g. AVFormatContext) so we keep as-is.
func GoStructName(cName string) string {
	return cName
}

// PublicTypeName converts a C type name to a clean public Go name.
// Strips "AV" prefix for public API.
//
// Examples:
//
//	AVFormatContext -> FormatContext
//	AVCodecContext  -> CodecContext
//	AVFrame         -> Frame
//	AVPacket        -> Packet
//	AVPixelFormat   -> PixelFormat
//	SwsContext      -> SwsContext (no AV prefix)
//	SwrContext      -> SwrContext (no AV prefix)
func PublicTypeName(cName string) string {
	s := cName
	// Strip AV prefix for public API
	if strings.HasPrefix(s, "AV") && len(s) > 2 && s[2] >= 'A' && s[2] <= 'Z' {
		s = s[2:]
	}
	return s
}

// GoFuncVarName converts a C function name to a Go variable name.
// e.g. avformat_open_input -> AvformatOpenInput
func GoFuncVarName(cName string) string {
	parts := strings.Split(cName, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if upper, ok := acronyms[strings.ToLower(p)]; ok {
			b.WriteString(upper)
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return b.String()
}

// PublicFuncName converts a C function name to a public Go function name.
// Keeps a short library prefix to avoid collisions between libs.
//
// Examples:
//
//	avformat_open_input    -> FormatOpenInput
//	avformat_free_context  -> FormatFreeContext
//	avcodec_free_context   -> CodecFreeContext
//	av_frame_alloc         -> FrameAlloc
//	sws_getContext         -> SwsGetContext
//	swr_alloc              -> SwrAlloc
//	avio_open              -> AvioOpen
func PublicFuncName(cName string) string {
	s := cName
	// Map library prefixes to short public prefixes
	prefixMap := []struct{ cPrefix, goPrefix string }{
		{"avformat_", "Format"},
		{"avcodec_", "Codec"},
		{"swscale_", "Swscale"},
		{"swresample_", "Swresample"},
		{"sws_", "Sws"},
		{"swr_", "Swr"},
		{"avio_", "Avio"},
		{"av_", ""},
	}
	goPrefix := ""
	for _, pm := range prefixMap {
		if strings.HasPrefix(s, pm.cPrefix) {
			goPrefix = pm.goPrefix
			s = s[len(pm.cPrefix):]
			break
		}
	}
	// PascalCase — preserve camelCase within segments
	parts := strings.Split(s, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if upper, ok := acronyms[strings.ToLower(p)]; ok {
			b.WriteString(upper)
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	return goPrefix + b.String()
}

// GoEnumName converts a C enum name to a Go type name.
// e.g. AVPixelFormat -> PixelFormat, AVCodecID -> CodecID
func GoEnumName(cName string) string {
	return PublicTypeName(cName)
}

// GoEnumValueName converts a C enum value to a Go constant name.
// e.g. AV_PIX_FMT_YUV420P -> PixFmtYUV420P
func GoEnumValueName(cName string) string {
	s := cName
	// Strip common AV_ prefix
	s = strings.TrimPrefix(s, "AV_")
	// Strip SWS_ prefix
	s = strings.TrimPrefix(s, "SWS_")

	// Convert UPPER_SNAKE to PascalCase
	parts := strings.Split(strings.ToLower(s), "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if upper, ok := acronyms[p]; ok {
			b.WriteString(upper)
		} else {
			b.WriteString(strings.ToUpper(p[:1]) + p[1:])
		}
	}
	result := b.String()
	if result == "" {
		return cName
	}
	return result
}
