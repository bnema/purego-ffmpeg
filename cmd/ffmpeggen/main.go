package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/emitter"
	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/model"
	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/overrides"
	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/parser"
)

type config struct {
	headersDir string
	outputDir  string
}

func main() {
	var cfg config
	flag.StringVar(&cfg.headersDir, "headers-dir", "/usr/include", "FFmpeg include root")
	flag.StringVar(&cfg.outputDir, "output-dir", ".", "project root output directory")
	flag.Parse()
	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	// Create all output directories
	dirs := []string{
		filepath.Join(cfg.outputDir, "internal/capi"),
		filepath.Join(cfg.outputDir, "internal/ports/out"),
		filepath.Join(cfg.outputDir, "internal/ports/out/mocks"),
		filepath.Join(cfg.outputDir, "internal/ports/in"),
		filepath.Join(cfg.outputDir, "ffmpeg"),
		filepath.Join(cfg.outputDir, "ffmpeg/mocks"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return fmt.Errorf("create dir %s: %w", d, err)
		}
	}

	// Step 1: Parse all headers needed by all domains
	headers, err := parseAllHeaders(cfg.headersDir)
	if err != nil {
		return err
	}

	// Step 2: Build domain viewmodels
	var domains []emitter.DomainData
	for _, domain := range overrides.Domains {
		dd := emitter.BuildDomainData(headers, domain)
		if len(dd.Functions) == 0 && len(dd.Accessors) == 0 {
			fmt.Printf("  SKIP %s (no functions or accessors matched)\n", domain.Name)
			continue
		}
		domains = append(domains, dd)
	}

	// Step 3: Emit per-domain files
	for _, dd := range domains {
		// internal/capi/{domain}_gen.go
		if err := emitFile(cfg.outputDir, filepath.Join("internal/capi", dd.Name+"_gen.go"), func() (string, error) {
			return emitter.EmitCAPIFile(dd)
		}); err != nil {
			return err
		}

		// internal/ports/out/{domain}_gen.go
		if err := emitFile(cfg.outputDir, filepath.Join("internal/ports/out", dd.Name+"_gen.go"), func() (string, error) {
			return emitter.EmitPortOut(dd)
		}); err != nil {
			return err
		}

		// internal/ports/out/mocks/{domain}_gen.go
		if err := emitFile(cfg.outputDir, filepath.Join("internal/ports/out/mocks", dd.Name+"_gen.go"), func() (string, error) {
			return emitter.EmitPortOutMock(dd)
		}); err != nil {
			return err
		}

		// internal/ports/in/{domain}_gen.go (only if domain has a PublicType)
		if dd.PublicType != "" {
			if err := emitFile(cfg.outputDir, filepath.Join("internal/ports/in", dd.Name+"_gen.go"), func() (string, error) {
				return emitter.EmitPortIn(dd)
			}); err != nil {
				return err
			}
		}

		// ffmpeg/{domain}_gen.go
		if err := emitFile(cfg.outputDir, filepath.Join("ffmpeg", dd.Name+"_gen.go"), func() (string, error) {
			return emitter.EmitPublicFileNew(dd)
		}); err != nil {
			return err
		}

		// ffmpeg/mocks/{domain}_gen.go (only if domain has a PublicType)
		if dd.PublicType != "" {
			if err := emitFile(cfg.outputDir, filepath.Join("ffmpeg/mocks", dd.Name+"_gen.go"), func() (string, error) {
				return emitter.EmitPublicMock(dd)
			}); err != nil {
				return err
			}
		}

		fmt.Printf("  %s: %d funcs, %d accessors\n", dd.Name, len(dd.Functions), len(dd.Accessors))
	}

	// Step 4: Emit cross-domain files

	// internal/capi/register_gen.go
	regData := emitter.BuildCAPIRegisterData(domains)
	if err := emitFile(cfg.outputDir, "internal/capi/register_gen.go", func() (string, error) {
		return emitter.EmitCAPIRegister(regData)
	}); err != nil {
		return err
	}

	// internal/capi/adapters_gen.go
	if err := emitFile(cfg.outputDir, "internal/capi/adapters_gen.go", func() (string, error) {
		return emitter.EmitCAPIAdapters(domains)
	}); err != nil {
		return err
	}

	// internal/ports/out/capi_gen.go
	if err := emitFile(cfg.outputDir, "internal/ports/out/capi_gen.go", func() (string, error) {
		return emitter.EmitPortOutCAPI(domains)
	}); err != nil {
		return err
	}

	// ffmpeg/types_gen.go
	typesData := emitter.BuildTypesData(headers, "ffmpeg", overrides.Enums, overrides.Structs)
	if err := emitFile(cfg.outputDir, "ffmpeg/types_gen.go", func() (string, error) {
		return emitter.EmitTypes(typesData)
	}); err != nil {
		return err
	}

	// internal/ports/in/types_gen.go
	inTypesData := emitter.BuildTypesData(headers, "in", overrides.Enums, overrides.Structs)
	if err := emitFile(cfg.outputDir, "internal/ports/in/types_gen.go", func() (string, error) {
		return emitter.EmitTypes(inTypesData)
	}); err != nil {
		return err
	}

	// ffmpeg/init_gen.go
	initData := emitter.BuildInitData(domains)
	if err := emitFile(cfg.outputDir, "ffmpeg/init_gen.go", func() (string, error) {
		return emitter.EmitInit(initData)
	}); err != nil {
		return err
	}

	fmt.Printf("\nGenerated %d domains, cross-domain files done.\n", len(domains))
	return nil
}

func emitFile(root, relPath string, emit func() (string, error)) error {
	code, err := emit()
	if err != nil {
		return fmt.Errorf("emit %s: %w", relPath, err)
	}
	fullPath := filepath.Join(root, relPath)
	if err := os.WriteFile(fullPath, []byte(code), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", relPath, err)
	}
	return nil
}

// parseAllHeaders parses all FFmpeg headers needed by the domain configuration.
// It builds a scope from all domains and parses each unique header file once.
func parseAllHeaders(headersDir string) ([]*model.Header, error) {
	// Collect all unique C function names across all domains
	allFuncs := make(map[string]bool)
	for _, d := range overrides.Domains {
		for _, fm := range d.Functions {
			allFuncs[fm.C] = true
		}
	}

	// Collect all enum names
	allEnums := make(map[string]bool)
	for _, e := range overrides.Enums {
		allEnums[e.C] = true
	}

	// Collect all struct names
	allStructs := make(map[string]bool)
	for _, s := range overrides.Structs {
		allStructs[s.C] = true
	}

	// Build the header file list — parse each library's main headers
	type headerInfo struct {
		path string
		lib  string
	}
	headerFiles := []headerInfo{
		// libavutil
		{"libavutil/error.h", "avutil"},
		{"libavutil/frame.h", "avutil"},
		{"libavutil/pixfmt.h", "avutil"},
		{"libavutil/samplefmt.h", "avutil"},
		{"libavutil/avutil.h", "avutil"},
		{"libavutil/rational.h", "avutil"},
		{"libavutil/mathematics.h", "avutil"},
		{"libavutil/dict.h", "avutil"},
		{"libavutil/log.h", "avutil"},
		{"libavutil/imgutils.h", "avutil"},
		{"libavutil/opt.h", "avutil"},
		{"libavutil/mem.h", "avutil"},
		// libavcodec
		{"libavcodec/avcodec.h", "avcodec"},
		{"libavcodec/codec.h", "avcodec"},
		{"libavcodec/codec_id.h", "avcodec"},
		{"libavcodec/packet.h", "avcodec"},
		// libavformat
		{"libavformat/avformat.h", "avformat"},
		{"libavformat/avio.h", "avformat"},
		// libavutil (hwcontext)
		{"libavutil/hwcontext.h", "avutil"},
		// libswscale
		{"libswscale/swscale.h", "swscale"},
		// libswresample
		{"libswresample/swresample.h", "swresample"},
		// libavfilter
		{"libavfilter/avfilter.h", "avfilter"},
		{"libavfilter/buffersrc.h", "avfilter"},
		{"libavfilter/buffersink.h", "avfilter"},
	}

	// Build a combined scope
	scope := &parser.Scope{
		Functions: allFuncs,
		Structs:   allStructs,
		Enums:     allEnums,
	}

	var headers []*model.Header
	for _, hf := range headerFiles {
		fullPath := filepath.Join(headersDir, hf.path)
		// Skip headers that don't exist on this system
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "  WARN: %s not found, skipping\n", hf.path)
			continue
		}
		header, err := parser.ParseFile(fullPath, hf.lib, scope)
		if err != nil {
			return nil, fmt.Errorf("parse %s: %w", hf.path, err)
		}
		headers = append(headers, header)
	}

	if len(headers) == 0 {
		return nil, fmt.Errorf("no FFmpeg headers found in %s; install FFmpeg development headers", headersDir)
	}

	return headers, nil
}
