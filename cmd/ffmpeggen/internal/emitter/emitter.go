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

// EmitRaw takes a parsed Header and returns formatted Go source for the raw layer.
func EmitRaw(header *model.Header) (string, error) {
	tmpl, err := template.New("raw").ParseFS(templateFS, "templates/raw_file.tmpl")
	if err != nil {
		return "", err
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
		return "", err
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
