package emitter

import (
	"strings"

	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/model"
	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/overrides"
)

// skipPublicFuncs lists C function names that should not be generated
// in the public layer because they have hand-written wrappers in support.go
// or their signatures cannot be correctly expressed in Go.
var skipPublicFuncs = map[string]bool{
	// av_strerror has an output buffer param (char*) that cannot work
	// with immutable Go strings. Handled by avError.Error() in support.go.
	"av_strerror": true,
}

// BuildPublicFileData converts parsed headers into the public API view model.
func BuildPublicFileData(header *model.Header) *PublicFileData {
	data := &PublicFileData{
		PackageName: "ffmpeg",
	}

	for i := range header.Structs {
		s := &header.Structs[i]
		w := buildWrapper(s)
		data.Wrappers = append(data.Wrappers, w)
	}

	for i := range header.Enums {
		data.Enums = append(data.Enums, buildEnum(&header.Enums[i]))
	}

	for i := range header.Functions {
		fn := &header.Functions[i]
		if skipPublicFuncs[fn.CName] {
			continue
		}
		data.FreeFunctions = append(data.FreeFunctions, buildFreeFunc(fn))
	}

	return data
}

func buildWrapper(s *model.Struct) WrapperData {
	w := WrapperData{
		Name:      s.PublicName,
		Doc:       s.Doc,
		RawGoName: s.GoName,
		IsOpaque:  s.IsOpaque,
	}

	for _, f := range s.Fields {
		w.Fields = append(w.Fields, WrapperFieldData{
			Name:       f.GoName,
			PublicType: resolvePublicFieldType(f.GoType),
			RawField:   f.GoName,
			Doc:        f.Doc,
		})
	}

	return w
}

func buildEnum(e *model.Enum) EnumData {
	ed := EnumData{
		Name: e.GoName,
		Doc:  e.Doc,
	}
	for _, v := range e.Values {
		ed.Values = append(ed.Values, EnumValueData{
			Name:  v.GoName,
			Value: v.Value,
		})
	}
	return ed
}

func buildFreeFunc(fn *model.Function) FreeFuncData {
	ff := FreeFuncData{
		Name:      model.PublicFuncName(fn.CName),
		Doc:       fn.Doc,
		RawGoName: fn.GoName,
	}

	for _, p := range fn.Params {
		pd := ParamData{
			Name:       p.GoName,
			PublicType: resolvePublicParamType(p.CType),
			CType:      p.CType,
			RawGoType:  p.GoType,
		}
		ct := normalizeConst(p.CType)
		pd.IsString = (ct == "char *" || ct == "char*")
		pd.IsPointer = p.GoType == "unsafe.Pointer"
		pd.IsDoublePointer = strings.Count(p.CType, "*") >= 2
		ff.Params = append(ff.Params, pd)
	}

	ret := strings.TrimSpace(fn.ReturnCType)
	if ret == "" || ret == "void" {
		ff.Return = ReturnData{IsVoid: true}
	} else {
		pubType := resolvePublicReturnType(ret)
		ff.Return = ReturnData{
			PublicType:     pubType,
			IsError:        isErrorReturn(ret, fn.CName),
			IsPointer:      fn.ReturnGoType == "unsafe.Pointer",
			IsStringReturn: IsStringReturn(fn),
		}
	}

	return ff
}

func normalizeConst(ctype string) string {
	ct := strings.ReplaceAll(ctype, "const ", "")
	return strings.TrimSpace(ct)
}

func resolvePublicFieldType(goType string) string {
	if goType == "unsafe.Pointer" {
		return "unsafe.Pointer"
	}
	return goType
}

func resolvePublicParamType(ctype string) string {
	ct := normalizeConst(ctype)
	if ct == "char *" || ct == "char*" {
		return "string"
	}
	// Default: use the raw Go type mapping
	return mapTypePublic(ctype)
}

func resolvePublicReturnType(ctype string) string {
	ct := normalizeConst(ctype)
	// const char* returns → string
	if ct == "char *" || ct == "char*" {
		return "string"
	}
	if strings.Contains(ct, "*") {
		return "unsafe.Pointer"
	}
	return mapTypePublic(ctype)
}

// IsStringReturn reports whether a function returns a C string (const char*).
func IsStringReturn(fn *model.Function) bool {
	ct := normalizeConst(strings.TrimSpace(fn.ReturnCType))
	return ct == "char *" || ct == "char*"
}

func mapTypePublic(ctype string) string {
	ct := normalizeConst(strings.TrimSpace(ctype))
	switch ct {
	case "int":
		return "int32"
	case "unsigned int", "unsigned":
		return "uint32"
	case "int64_t":
		return "int64"
	case "uint64_t":
		return "uint64"
	case "int32_t":
		return "int32"
	case "uint32_t":
		return "uint32"
	case "int8_t":
		return "int8"
	case "uint8_t":
		return "uint8"
	case "int16_t":
		return "int16"
	case "uint16_t":
		return "uint16"
	case "long":
		return "int64"
	case "unsigned long":
		return "uint64"
	case "size_t":
		return "uintptr"
	case "float":
		return "float32"
	case "double":
		return "float64"
	case "void":
		return ""
	case "AVRational":
		return "AVRational"
	case "AVChannelLayout":
		return "AVChannelLayout"
	}
	if strings.HasPrefix(ct, "enum ") {
		return "int32"
	}
	if strings.Contains(ct, "*") {
		return "unsafe.Pointer"
	}
	return "int32"
}

// isErrorReturn returns true if the function returns int and follows FFmpeg convention
// of returning negative AVERROR codes on failure.
func isErrorReturn(retType string, funcName string) bool {
	ct := normalizeConst(strings.TrimSpace(retType))
	if ct != "int" {
		return false
	}
	// Most FFmpeg functions returning int use AVERROR convention.
	// Exceptions: functions that return counts, booleans, etc.
	noErrorFuncs := map[string]bool{
		"avcodec_is_open":      true,
		"av_frame_is_writable": true,
	}
	return !noErrorFuncs[funcName]
}

// titleCase returns s with the first character uppercased.
// Used instead of the deprecated strings.Title.
func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// BuildDomainData takes parsed model functions and an overrides domain config and
// produces a DomainData viewmodel ready for template execution.
func BuildDomainData(headers []*model.Header, domain overrides.Domain) DomainData {
	dd := DomainData{
		Name:          domain.Name,
		PortInterface: domain.PortInterface,
		PublicType:    domain.PublicType,
		Library:       domain.Library,
	}

	// Build a lookup map of all parsed functions across all headers for the domain's library.
	funcMap := make(map[string]*model.Function)
	for _, h := range headers {
		for i := range h.Functions {
			funcMap[h.Functions[i].CName] = &h.Functions[i]
		}
	}

	// Match domain FuncMap entries to parsed functions.
	for _, fm := range domain.Functions {
		fn, ok := funcMap[fm.C]
		if !ok {
			// Function not found in parsed headers — skip silently.
			// This can happen if the header doesn't exist on this system.
			continue
		}
		dfd := DomainFuncData{
			CName:      fm.C,
			GoMethod:   fm.Go,
			RawVarName: fn.CName, // keep snake_case for purego var name
		}

		// Build params from parsed function.
		for _, p := range fn.Params {
			dfd.Params = append(dfd.Params, RawParamData{
				Name:   p.GoName,
				GoType: p.GoType,
			})
		}

		// Build return type.
		ret := strings.TrimSpace(fn.ReturnCType)
		if ret == "" || ret == "void" {
			dfd.Return = RawReturnData{IsVoid: true}
		} else {
			dfd.Return = RawReturnData{
				GoType:  fn.ReturnGoType,
				IsError: isErrorReturn(ret, fn.CName),
			}
		}

		dd.Functions = append(dd.Functions, dfd)
	}

	// Build accessors from overrides.
	for _, acc := range domain.Accessors {
		dd.Accessors = append(dd.Accessors, AccessorData{
			Struct:    acc.Struct,
			Field:     acc.Field,
			GoName:    acc.GoName,
			GoType:    acc.Type,
			Offset:    acc.Offset,
			HasSetter: true, // all overrides-defined accessors are mutable for now
		})
	}

	return dd
}

// BuildTypesData builds TypesData from parsed enum values and overrides enum/struct definitions.
func BuildTypesData(headers []*model.Header, packageName string, enumDefs []overrides.EnumDef, structDefs []overrides.StructDef) TypesData {
	td := TypesData{
		PackageName: packageName,
	}

	// Build enum lookup from all parsed headers.
	enumMap := make(map[string]*model.Enum)
	for _, h := range headers {
		for i := range h.Enums {
			enumMap[h.Enums[i].CName] = &h.Enums[i]
		}
	}

	// Match overrides enum definitions to parsed enums.
	for _, ed := range enumDefs {
		parsedEnum, ok := enumMap[ed.C]
		if !ok {
			// Enum not found — create empty type definition.
			td.Enums = append(td.Enums, EnumData{Name: ed.Go})
			continue
		}
		td.Enums = append(td.Enums, buildEnum(parsedEnum))
	}

	// Build value structs from overrides.
	for _, sd := range structDefs {
		vs := ValueStructData{GoName: sd.Go}
		for _, f := range sd.Fields {
			vs.Fields = append(vs.Fields, ValueStructFieldData{
				Name: f.Name,
				Type: f.Type,
			})
		}
		td.Structs = append(td.Structs, vs)
	}

	return td
}

// BuildInitData creates an InitData viewmodel from a slice of domain viewmodels.
func BuildInitData(domains []DomainData) InitData {
	return InitData{Domains: domains}
}

// BuildCAPIRegisterData groups domains by library and creates registration entries.
func BuildCAPIRegisterData(domains []DomainData) CAPIRegisterData {
	type libEntry struct {
		handleField   string
		registerFuncs []string
	}

	libOrder := []string{} // preserve order
	libMap := make(map[string]*libEntry)

	for _, d := range domains {
		handleField := libraryToHandleField(d.Library)
		if _, ok := libMap[handleField]; !ok {
			libMap[handleField] = &libEntry{handleField: handleField}
			libOrder = append(libOrder, handleField)
		}
		registerFunc := "Register" + titleCase(d.Name)
		libMap[handleField].registerFuncs = append(libMap[handleField].registerFuncs, registerFunc)
	}

	var data CAPIRegisterData
	for _, hf := range libOrder {
		le := libMap[hf]
		data.Libraries = append(data.Libraries, CAPILibraryData{
			HandleField:   le.handleField,
			RegisterFuncs: le.registerFuncs,
		})
	}
	return data
}

// libraryToHandleField maps a library name to its handle field name in the CAPI struct.
func libraryToHandleField(lib string) string {
	switch lib {
	case "libavutil":
		return "Avutil"
	case "libavcodec":
		return "Avcodec"
	case "libavformat":
		return "Avformat"
	case "libswscale":
		return "Swscale"
	case "libswresample":
		return "Swresample"
	default:
		return lib
	}
}
