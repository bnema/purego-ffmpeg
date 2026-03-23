package parser

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/model"
)

// Scope defines which symbols to extract from headers.
type Scope struct {
	Functions map[string]bool // C function names to include
	Structs   map[string]bool // C struct names to include
	Enums     map[string]bool // C enum names to include
}

var (
	// Match typedef struct Name { ... } Name;
	typedefStructRE = regexp.MustCompile(`(?s)typedef\s+struct\s+(\w+)\s*\{(.*?)\}\s*(\w+)\s*;`)
	// Match simple typedef struct Name Name; (opaque forward declaration)
	opaqueTypedefRE = regexp.MustCompile(`typedef\s+struct\s+(\w+)\s+(\w+)\s*;`)
	// Match enum Name { ... };
	enumRE = regexp.MustCompile(`(?s)enum\s+(\w+)\s*\{(.*?)\}\s*;`)
	// Match function declarations: type name(params);
	// This is broad - we filter by scope. Handles multi-word return types.
	funcRE = regexp.MustCompile(`^(\w[\w\s\*]+?)\s+(\*?\w+)\s*\(([^)]*)\)\s*;`)
	// Match __attribute__((...)) patterns
	attrRE = regexp.MustCompile(`__attribute__\s*\(\([^)]*\)\)`)
)

// ParseFile reads an FFmpeg header file and returns a parsed Header.
func ParseFile(path string, lib string, scope *Scope) (*model.Header, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(path, lib, data, scope)
}

// Parse parses the given FFmpeg header content.
func Parse(path string, lib string, data []byte, scope *Scope) (*model.Header, error) {
	rawSource := string(data)
	docIdx := buildDocIndex(rawSource)

	clean := stripComments(data)
	clean = joinLines(clean)

	out := &model.Header{Path: path, Lib: lib}

	// Parse typedef structs
	for _, match := range typedefStructRE.FindAllSubmatch(clean, -1) {
		name := string(match[3])
		if scope != nil && !scope.Structs[name] {
			continue
		}
		st := parseStruct(string(match[2]), name)
		populateStructDoc(&st, rawSource, docIdx)
		out.Structs = append(out.Structs, st)
	}

	// Parse opaque typedef structs
	for _, match := range opaqueTypedefRE.FindAllSubmatch(clean, -1) {
		name := string(match[2])
		if scope != nil && !scope.Structs[name] {
			continue
		}
		// Skip if we already parsed the full struct
		found := false
		for _, s := range out.Structs {
			if s.CName == name {
				found = true
				break
			}
		}
		if found {
			continue
		}
		st := model.Struct{
			CName:      name,
			GoName:     model.GoStructName(name),
			PublicName: model.PublicTypeName(name),
			Kind:       "opaque",
			IsOpaque:   true,
		}
		out.Structs = append(out.Structs, st)
	}

	// Parse enums
	for _, match := range enumRE.FindAllSubmatch(clean, -1) {
		name := string(match[1])
		if scope != nil && !scope.Enums[name] {
			continue
		}
		e := parseEnum(name, string(match[2]))
		populateEnumDoc(&e, rawSource, docIdx)
		out.Enums = append(out.Enums, e)
	}

	// Parse functions - line by line from cleaned source
	for _, line := range bytes.Split(clean, []byte("\n")) {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		match := funcRE.FindSubmatch(trimmed)
		if match == nil {
			continue
		}
		retType := strings.TrimSpace(string(match[1]))
		name := string(match[2])
		// Handle pointer return: "AVFrame *av_frame_alloc"
		if strings.HasPrefix(name, "*") {
			retType += " *"
			name = name[1:]
		}
		params := string(match[3])

		if scope != nil && !scope.Functions[name] {
			continue
		}

		fn := parseFunction(name, retType, params)
		populateFuncDoc(&fn, rawSource, docIdx)
		out.Functions = append(out.Functions, fn)
	}

	return out, nil
}

func parseStruct(body, name string) model.Struct {
	st := model.Struct{
		CName:      name,
		GoName:     model.GoStructName(name),
		PublicName: model.PublicTypeName(name),
		Kind:       "data",
	}

	for _, raw := range strings.Split(body, ";") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		// Skip preprocessor directives and nested struct/union/enum definitions.
		// Note: "enum Foo bar" (no braces) is an enum-typed field, not a declaration.
		if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "struct ") ||
			strings.HasPrefix(line, "union ") ||
			(strings.HasPrefix(line, "enum ") && strings.Contains(line, "{")) ||
			strings.Contains(line, "{") || strings.Contains(line, "}") {
			continue
		}

		// Handle fixed-size arrays: type name[size]
		if idx := strings.Index(line, "["); idx >= 0 {
			endIdx := strings.Index(line, "]")
			if endIdx > idx {
				arraySize := strings.TrimSpace(line[idx+1 : endIdx])
				line = line[:idx] // Remove the [size] part
				parts := strings.Fields(line)
				if len(parts) < 2 {
					continue
				}
				fieldName := parts[len(parts)-1]
				fieldType := strings.Join(parts[:len(parts)-1], " ")
				st.Fields = append(st.Fields, model.Field{
					CName:     fieldName,
					GoName:    goFieldName(fieldName),
					CType:     fieldType,
					GoType:    "[" + arraySize + "]" + mapType(fieldType),
					IsArray:   true,
					ArraySize: arraySize,
				})
				continue
			}
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		fieldName := parts[len(parts)-1]
		fieldType := strings.Join(parts[:len(parts)-1], " ")

		// Handle pointer in field name
		stars := ""
		for strings.HasPrefix(fieldName, "*") {
			stars += "*"
			fieldName = fieldName[1:]
		}
		if stars != "" {
			fieldType += " " + stars
		}

		st.Fields = append(st.Fields, model.Field{
			CName:  fieldName,
			GoName: goFieldName(fieldName),
			CType:  fieldType,
			GoType: mapType(fieldType),
		})
	}
	return st
}

func parseFunction(name, ret, params string) model.Function {
	return model.Function{
		CName:        name,
		GoName:       model.GoFuncVarName(name),
		Params:       parseParams(params),
		ReturnCType:  strings.TrimSpace(ret),
		ReturnGoType: mapType(ret),
	}
}

func parseParams(raw string) []model.Param {
	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "void" {
		return nil
	}
	parts := splitParams(raw)
	params := make([]model.Param, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" || part == "void" || part == "..." {
			continue
		}
		fields := strings.Fields(part)
		if len(fields) == 0 {
			continue
		}
		name := fields[len(fields)-1]
		// Strip trailing [] (C array parameter syntax, treated as pointer)
		if strings.HasSuffix(name, "[]") {
			name = strings.TrimSuffix(name, "[]")
		}
		stars := ""
		for strings.HasPrefix(name, "*") {
			stars += "*"
			name = name[1:]
		}
		ctype := strings.TrimSpace(strings.Join(fields[:len(fields)-1], " ")) + stars
		params = append(params, model.Param{
			CName:   name,
			GoName:  goParamName(name),
			CType:   ctype,
			GoType:  mapType(ctype),
			IsConst: strings.Contains(ctype, "const "),
			Pointer: strings.Count(ctype, "*"),
		})
	}
	return params
}

func splitParams(raw string) []string {
	var parts []string
	depth := 0
	start := 0
	for i, ch := range raw {
		switch ch {
		case '(':
			depth++
		case ')':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, raw[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, raw[start:])
	return parts
}

func parseEnum(name, body string) model.Enum {
	result := model.Enum{
		CName:  name,
		GoName: model.GoEnumName(name),
	}

	nextVal := 0
	nameToVal := map[string]int{}
	seen := map[string]bool{}

	for _, line := range strings.Split(body, ",") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Strip inline comments
		if idx := strings.Index(line, "//"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if idx := strings.Index(line, "/*"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		cname := strings.TrimSpace(parts[0])
		if cname == "" {
			continue
		}
		if seen[cname] {
			continue
		}
		seen[cname] = true

		val := ""
		if len(parts) == 2 {
			val = strings.TrimSpace(parts[1])
		}

		if val == "" {
			val = fmt.Sprintf("%d", nextVal)
			nameToVal[cname] = nextVal
			nextVal++
		} else if n, err := strconv.Atoi(val); err == nil {
			nameToVal[cname] = n
			nextVal = n + 1
		} else if resolved, ok := nameToVal[val]; ok {
			val = fmt.Sprintf("%d", resolved)
			nameToVal[cname] = resolved
			nextVal = resolved + 1
		} else {
			// Complex expression - keep as-is
			nameToVal[cname] = nextVal
			nextVal++
		}

		result.Values = append(result.Values, model.EnumValue{
			CName:  cname,
			GoName: model.GoEnumValueName(cname),
			Value:  val,
		})
	}
	return result
}

// mapType maps a C type string to a Go type string.
func mapType(ctype string) string {
	ctype = strings.TrimSpace(ctype)
	ctype = strings.Join(strings.Fields(ctype), " ")

	// Remove const qualifier for mapping
	ct := strings.ReplaceAll(ctype, "const ", "")
	ct = strings.TrimSpace(ct)

	isPtr := strings.Contains(ct, "*")

	// void* -> unsafe.Pointer
	if ct == "void *" || ct == "void*" {
		return "unsafe.Pointer"
	}
	if ct == "void **" || ct == "void**" {
		return "unsafe.Pointer"
	}

	// void (no pointer) -> empty
	if ct == "void" {
		return ""
	}

	// int -> int32
	if ct == "int" {
		return "int32"
	}

	// unsigned int / unsigned -> uint32
	if ct == "unsigned int" || ct == "unsigned" {
		return "uint32"
	}

	// size_t -> uintptr
	if ct == "size_t" {
		return "uintptr"
	}

	// stdint types
	switch ct {
	case "int8_t":
		return "int8"
	case "int16_t":
		return "int16"
	case "int32_t":
		return "int32"
	case "int64_t":
		return "int64"
	case "uint8_t":
		return "uint8"
	case "uint16_t":
		return "uint16"
	case "uint32_t":
		return "uint32"
	case "uint64_t":
		return "uint64"
	}

	// float/double
	if ct == "float" {
		return "float32"
	}
	if ct == "double" {
		return "float64"
	}

	// char* / const char* -> *byte
	if ct == "char *" || ct == "char*" {
		return "*byte"
	}

	// char** -> unsafe.Pointer
	if ct == "char **" || ct == "char**" {
		return "unsafe.Pointer"
	}

	// uint8_t* -> unsafe.Pointer
	if ct == "uint8_t *" || ct == "uint8_t*" {
		return "unsafe.Pointer"
	}

	// enum types without pointer -> int32
	if strings.HasPrefix(ct, "enum ") && !isPtr {
		return "int32"
	}

	// Any pointer type -> unsafe.Pointer
	if isPtr {
		return "unsafe.Pointer"
	}

	// AVRational (value type) -> keep as struct name
	if ct == "AVRational" {
		return "AVRational"
	}

	// AVChannelLayout (value type)
	if ct == "AVChannelLayout" {
		return "AVChannelLayout"
	}

	// Known struct-by-value types — mapped to fixed-size byte arrays
	// to preserve correct struct layout when embedded in other structs.
	// Sizes are for 64-bit Linux (FFmpeg 7.x/8.x).
	knownStructSizes := map[string]int{
		"AVPacket": 104,
	}
	if size, ok := knownStructSizes[ct]; ok {
		return fmt.Sprintf("[%d]byte", size)
	}

	// Unknown non-pointer
	return "uintptr"
}

// goFieldName converts a C field name to a Go exported field name.
func goFieldName(cname string) string {
	parts := strings.Split(cname, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		low := strings.ToLower(p)
		switch low {
		case "id":
			b.WriteString("ID")
		case "url":
			b.WriteString("URL")
		case "io":
			b.WriteString("IO")
		case "pts":
			b.WriteString("PTS")
		case "dts":
			b.WriteString("DTS")
		case "fps":
			b.WriteString("FPS")
		case "nb":
			b.WriteString("Nb")
		case "hw":
			b.WriteString("HW")
		default:
			b.WriteString(strings.ToUpper(low[:1]) + low[1:])
		}
	}
	return b.String()
}

// goParamName converts a C parameter name to a Go parameter name (lowercase first letter).
func goParamName(cname string) string {
	name := goFieldName(cname)
	if len(name) == 0 {
		return cname
	}
	// Lowercase first letter
	name = strings.ToLower(name[:1]) + name[1:]
	// Escape Go keywords
	if goKeywords[name] {
		name = name + "_"
	}
	return name
}

// goKeywords are Go reserved words.
var goKeywords = map[string]bool{
	"break": true, "case": true, "chan": true, "const": true, "continue": true,
	"default": true, "defer": true, "else": true, "fallthrough": true,
	"for": true, "func": true, "go": true, "goto": true, "if": true,
	"import": true, "interface": true, "map": true, "package": true,
	"range": true, "return": true, "select": true, "struct": true,
	"switch": true, "type": true, "var": true,
}

// stripComments removes C comments and preprocessor directives.
func stripComments(data []byte) []byte {
	// First strip block comments
	for {
		start := bytes.Index(data, []byte("/*"))
		if start < 0 {
			break
		}
		end := bytes.Index(data[start+2:], []byte("*/"))
		if end < 0 {
			break
		}
		data = append(data[:start], data[start+2+end+2:]...)
	}

	// Strip FFmpeg attribute macros
	ffmpegMacros := []string{
		"av_warn_unused_result",
		"av_nonnull",
		"attribute_deprecated",
		"av_pure",
		"av_const",
		"av_cold",
		"av_noreturn",
	}
	for _, macro := range ffmpegMacros {
		data = bytes.ReplaceAll(data, []byte(macro), nil)
	}
	// Strip __attribute__((...)) patterns
	data = attrRE.ReplaceAll(data, nil)

	lines := bytes.Split(data, []byte("\n"))
	out := make([][]byte, 0, len(lines))
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		if bytes.HasPrefix(trimmed, []byte("//")) {
			continue
		}
		if bytes.HasPrefix(trimmed, []byte("#")) {
			continue
		}
		// Strip inline comments
		if idx := bytes.Index(line, []byte("//")); idx >= 0 {
			line = bytes.TrimRightFunc(line[:idx], unicode.IsSpace)
			if len(bytes.TrimSpace(line)) == 0 {
				continue
			}
		}
		out = append(out, line)
	}
	return bytes.Join(out, []byte("\n"))
}

// joinLines collapses multi-line declarations into single lines.
func joinLines(data []byte) []byte {
	lines := bytes.Split(data, []byte("\n"))
	var result [][]byte
	var buf []byte
	for _, line := range lines {
		trimmed := bytes.TrimSpace(line)
		if len(trimmed) == 0 {
			if buf != nil {
				result = append(result, buf)
				buf = nil
			}
			continue
		}
		if buf == nil {
			buf = bytes.TrimSpace(line)
		} else {
			buf = append(buf, ' ')
			buf = append(buf, bytes.TrimSpace(line)...)
		}
		last := trimmed[len(trimmed)-1]
		if last == ';' || last == '{' || last == '}' {
			result = append(result, buf)
			buf = nil
		}
	}
	if buf != nil {
		result = append(result, buf)
	}
	return bytes.Join(result, []byte("\n"))
}

// --- Doc extraction ---

type docBlock struct {
	startLine int
	endLine   int
	lines     []string
}

type docIndex struct {
	blocks []docBlock
}

func buildDocIndex(source string) *docIndex {
	lines := strings.Split(source, "\n")
	idx := &docIndex{}
	var current *docBlock

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "/**") {
			current = &docBlock{startLine: i}
			current.lines = append(current.lines, trimmed)
			if strings.Contains(trimmed, "*/") {
				current.endLine = i
				idx.blocks = append(idx.blocks, *current)
				current = nil
			}
		} else if current != nil {
			current.lines = append(current.lines, trimmed)
			if strings.Contains(trimmed, "*/") {
				current.endLine = i
				idx.blocks = append(idx.blocks, *current)
				current = nil
			}
		}
	}
	return idx
}

func (idx *docIndex) forLine(line int) *docBlock {
	for i := len(idx.blocks) - 1; i >= 0; i-- {
		b := &idx.blocks[i]
		if b.endLine < line && line-b.endLine <= 3 {
			return b
		}
	}
	return nil
}

func findLineOf(source, needle string) int {
	lines := strings.Split(source, "\n")
	for i, line := range lines {
		if strings.Contains(line, needle) {
			return i
		}
	}
	return -1
}

func cleanDoc(lines []string) string {
	var cleaned []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		l = strings.TrimPrefix(l, "/**")
		l = strings.TrimPrefix(l, "* ")
		l = strings.TrimPrefix(l, "*")
		l = strings.TrimSuffix(l, "*/")
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "@") {
			continue
		}
		cleaned = append(cleaned, l)
	}
	if len(cleaned) == 0 {
		return ""
	}
	return strings.Join(cleaned, " ")
}

func populateStructDoc(st *model.Struct, rawSource string, docIdx *docIndex) {
	needle := "typedef struct " + st.CName
	line := findLineOf(rawSource, needle)
	if line >= 0 {
		if db := docIdx.forLine(line); db != nil {
			st.Doc = cleanDoc(db.lines)
		}
	}
}

func populateEnumDoc(e *model.Enum, rawSource string, docIdx *docIndex) {
	needle := "enum " + e.CName
	line := findLineOf(rawSource, needle)
	if line >= 0 {
		if db := docIdx.forLine(line); db != nil {
			e.Doc = cleanDoc(db.lines)
		}
	}
}

func populateFuncDoc(fn *model.Function, rawSource string, docIdx *docIndex) {
	needle := fn.CName + "("
	line := findLineOf(rawSource, needle)
	if line >= 0 {
		if db := docIdx.forLine(line); db != nil {
			fn.Doc = cleanDoc(db.lines)
		}
	}
}
