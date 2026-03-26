package emitter

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

// === Hexagonal viewmodel types ===

// DomainData is the top-level viewmodel for one domain (e.g., "format", "codec").
// Used to generate all layers for that domain.
type DomainData struct {
	Name          string // "format"
	PortInterface string // "FormatCAPI"
	PublicType    string // "FormatContext" (empty for free-func-only domains like "util")
	Library       string // "libavformat"
	Functions     []DomainFuncData
	Accessors     []AccessorData
	FreeMethod    string // CAPI method name to call in Free() — e.g., "FreeContext", "FreePtr"
	AllocMethod   string // CAPI method name for zero-arg allocation — e.g., "Alloc", "AllocContext"
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
