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
