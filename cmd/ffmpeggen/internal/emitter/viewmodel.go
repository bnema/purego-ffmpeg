package emitter

// PublicFileData holds all data needed to render the public API file.
type PublicFileData struct {
	PackageName   string
	Wrappers      []WrapperData
	Enums         []EnumData
	FreeFunctions []FreeFuncData
}

// NeedsUnsafe reports whether the generated public file needs "unsafe" imported.
func (d *PublicFileData) NeedsUnsafe() bool {
	if len(d.Wrappers) > 0 {
		return true
	}
	for _, ff := range d.FreeFunctions {
		for _, p := range ff.Params {
			if p.IsPointer || p.IsDoublePointer || p.PublicType == "unsafe.Pointer" {
				return true
			}
		}
		if ff.Return.IsPointer || ff.Return.PublicType == "unsafe.Pointer" {
			return true
		}
	}
	return false
}

// NeedsRaw reports whether the generated public file needs the raw package imported.
func (d *PublicFileData) NeedsRaw() bool {
	return len(d.Wrappers) > 0 || len(d.FreeFunctions) > 0
}

// NeedsRuntime reports whether the generated public file needs "runtime" imported.
func (d *PublicFileData) NeedsRuntime() bool {
	for _, w := range d.Wrappers {
		if w.HasClose {
			return true
		}
	}
	for _, ff := range d.FreeFunctions {
		for _, p := range ff.Params {
			if p.IsString {
				return true
			}
		}
	}
	return false
}

// WrapperData represents a public wrapper for an FFmpeg struct.
type WrapperData struct {
	Name      string // "FormatContext"
	Doc       string
	RawGoName string // "AVFormatContext"
	IsOpaque  bool   // true if struct has no exposed fields
	HasClose  bool   // true if there's a known free function
	CloseFunc string // raw function name for cleanup, e.g., "AvformatCloseInput"
	Fields    []WrapperFieldData
}

// WrapperFieldData represents a single accessor field on a public wrapper.
type WrapperFieldData struct {
	Name       string // "NbStreams"
	PublicType string // "int32"
	RawField   string // "NbStreams" (raw struct field name)
	Doc        string
}

// EnumData represents a generated enum type.
type EnumData struct {
	Name     string // "PixelFormat"
	Doc      string
	Unsigned bool
	Values   []EnumValueData
}

// EnumValueData represents a single enum constant.
type EnumValueData struct {
	Name  string // "PixFmtYuv420p"
	Value string // "0"
}

// FreeFuncData represents a wrapped free function.
type FreeFuncData struct {
	Name      string // "OpenInput"
	Doc       string
	RawGoName string // "AvformatOpenInput"
	Params    []ParamData
	Return    ReturnData
}

// ParamData represents a single function parameter.
type ParamData struct {
	Name            string
	PublicType      string
	CType           string
	IsString        bool // const char* -> string
	IsPointer       bool // unsafe.Pointer
	IsDoublePointer bool // e.g., AVFormatContext**
	IsEnum          bool
	RawGoType       string // raw Go type for direct pass-through
}

// ReturnData describes the return type of a wrapped function.
type ReturnData struct {
	PublicType     string
	IsVoid         bool
	IsError        bool // int return that represents AVERROR
	IsPointer      bool
	IsEnum         bool
	IsStringReturn bool // const char* return → Go string
}

// === New hexagonal viewmodel types ===

// DomainData is the top-level viewmodel for one domain (e.g., "format", "codec").
// Used to generate all layers for that domain.
type DomainData struct {
	Name          string // "format"
	PortInterface string // "FormatCAPI"
	PublicType    string // "FormatContext" (empty for free-func-only domains like "util")
	Library       string // "libavformat"
	Functions     []DomainFuncData
	Accessors     []AccessorData
}

// DomainFuncData represents one C→Go function mapping within a domain.
type DomainFuncData struct {
	CName      string // "avformat_open_input"
	GoMethod   string // "OpenInput"
	RawVarName string // "avformat_open_input" (purego var name, snake_case)
	Params     []RawParamData
	Return     RawReturnData
}

// RawParamData is a parameter in the raw C signature.
type RawParamData struct {
	Name   string // "ctx"
	GoType string // "unsafe.Pointer", "*unsafe.Pointer", "*byte", "int32", etc.
}

// RawReturnData is the return type in the raw C signature.
type RawReturnData struct {
	GoType  string // "int32", "unsafe.Pointer", ""
	IsVoid  bool
	IsError bool // int32 return that represents AVERROR
}

// AccessorData represents an offset-based struct field accessor.
type AccessorData struct {
	Struct    string // "AVCodecContext"
	Field     string // "width"
	GoName    string // "Width"
	GoType    string // "int32"
	Offset    int
	HasSetter bool // true for mutable fields
}

// TypesData holds all enums and value structs for the types_gen.go file.
type TypesData struct {
	PackageName string
	Enums       []EnumData // reuse existing EnumData
	Structs     []ValueStructData
}

// ValueStructData represents a simple C value struct (e.g., AVRational).
type ValueStructData struct {
	GoName string
	Fields []ValueStructFieldData
}

// ValueStructFieldData represents a field in a value struct.
type ValueStructFieldData struct {
	Name string
	Type string
}

// InitData holds data for the composition root init_gen.go template.
type InitData struct {
	Domains []DomainData // all domains for building accessor funcs
}

// CAPIRegisterData holds data for the register_gen.go aggregator.
type CAPIRegisterData struct {
	Libraries []CAPILibraryData
}

// CAPILibraryData maps a library to its registration function names.
type CAPILibraryData struct {
	HandleField   string   // "Avformat"
	RegisterFuncs []string // ["RegisterAvformat", "RegisterAvio"]
}
