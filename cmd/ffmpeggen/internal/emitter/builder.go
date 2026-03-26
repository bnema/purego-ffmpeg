package emitter

import (
	"fmt"
	"os"
	"strings"

	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/model"
	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/overrides"
)

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

func normalizeConst(ctype string) string {
	ct := strings.ReplaceAll(ctype, "const ", "")
	return strings.TrimSpace(ct)
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
			fmt.Fprintf(os.Stderr, "  WARN: %s not found in parsed headers for domain %s\n", fm.C, domain.Name)
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

	// Find the free method for this domain.
	for _, f := range dd.Functions {
		if strings.HasPrefix(f.GoMethod, "Free") {
			dd.FreeMethod = f.GoMethod
			break
		}
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
