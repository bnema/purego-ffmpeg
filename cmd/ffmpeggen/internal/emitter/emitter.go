package emitter

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/model"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// hexFuncMap is the shared FuncMap used by all hexagonal layer templates.
var hexFuncMap = template.FuncMap{
	"title": func(s string) string {
		if s == "" {
			return s
		}
		return strings.ToUpper(s[:1]) + s[1:]
	},
	"lower": func(s string) string {
		if s == "" {
			return s
		}
		return strings.ToLower(s[:1]) + s[1:]
	},
	// trimSuffix receives the pipeline value as the last argument:
	// {{.Name | trimSuffix "CAPI"}} → trimSuffix("CAPI", .Name)
	"trimSuffix": func(suffix, s string) string {
		return strings.TrimSuffix(s, suffix)
	},
	"lowerFirst": func(s string) string {
		if s == "" {
			return s
		}
		return strings.ToLower(s[:1]) + s[1:]
	},
	"upper": strings.ToUpper,
}

// renderTemplate parses a single template file from the embedded FS, executes it
// with the given data, and returns gofmt-formatted Go source.
func renderTemplate(tmplName string, data interface{}) (string, error) {
	tmpl, err := template.New(tmplName).Funcs(hexFuncMap).ParseFS(templateFS, "templates/"+tmplName)
	if err != nil {
		return "", fmt.Errorf("parse %s: %w", tmplName, err)
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, tmplName, data); err != nil {
		return "", fmt.Errorf("execute %s: %w", tmplName, err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format %s: %w\n---\n%s", tmplName, err, buf.String())
	}
	return string(formatted), nil
}

// EmitRaw takes a parsed Header and returns formatted Go source for the raw layer.
func EmitRaw(header *model.Header) (string, error) {
	tmpl, err := template.New("raw").ParseFS(templateFS, "templates/raw_file.tmpl")
	if err != nil {
		return "", fmt.Errorf("parse raw templates: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "raw_file.tmpl", header); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}

// EmitPublic takes a PublicFileData view model and returns formatted Go source.
func EmitPublic(data *PublicFileData) (string, error) {
	funcMap := template.FuncMap{
		"lower": func(s string) string {
			if s == "" {
				return s
			}
			return strings.ToLower(s[:1]) + s[1:]
		},
	}

	tmpl, err := template.New("public").Funcs(funcMap).ParseFS(templateFS,
		"templates/public_file.tmpl",
		"templates/wrapper.tmpl",
		"templates/enums.tmpl",
		"templates/free_func.tmpl",
	)
	if err != nil {
		return "", fmt.Errorf("parse public templates: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "public_file.tmpl", data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return "", fmt.Errorf("format source: %w\n%s", err, buf.String())
	}
	return string(formatted), nil
}

// EmitCAPIFile generates internal/capi/{domain}_gen.go
func EmitCAPIFile(domain DomainData) (string, error) {
	return renderTemplate("capi_file.tmpl", domain)
}

// EmitCAPIRegister generates internal/capi/register_gen.go
func EmitCAPIRegister(data CAPIRegisterData) (string, error) {
	return renderTemplate("capi_register.tmpl", data)
}

// EmitCAPIAdapters generates internal/capi/adapters_gen.go.
// Takes a wrapper struct since templates need .Domains field access.
func EmitCAPIAdapters(domains []DomainData) (string, error) {
	return renderTemplate("capi_adapters.tmpl", struct{ Domains []DomainData }{domains})
}

// EmitPortOut generates internal/ports/out/{domain}_gen.go
func EmitPortOut(domain DomainData) (string, error) {
	return renderTemplate("port_out.tmpl", domain)
}

// EmitPortOutMock generates internal/ports/out/mocks/{domain}_gen.go
func EmitPortOutMock(domain DomainData) (string, error) {
	return renderTemplate("port_out_mock.tmpl", domain)
}

// EmitPortIn generates internal/ports/in/{domain}_gen.go
func EmitPortIn(domain DomainData) (string, error) {
	return renderTemplate("port_in.tmpl", domain)
}

// EmitPublicFileNew generates ffmpeg/{domain}_gen.go
func EmitPublicFileNew(domain DomainData) (string, error) {
	return renderTemplate("public_file_new.tmpl", domain)
}

// EmitPublicMock generates ffmpeg/mocks/{domain}_gen.go
func EmitPublicMock(domain DomainData) (string, error) {
	return renderTemplate("public_mock.tmpl", domain)
}

// EmitTypes generates types_gen.go (used for both ffmpeg/ and internal/ports/in/)
func EmitTypes(data TypesData) (string, error) {
	return renderTemplate("types_new.tmpl", data)
}

// EmitInit generates ffmpeg/init_gen.go
func EmitInit(data InitData) (string, error) {
	return renderTemplate("init.tmpl", data)
}
