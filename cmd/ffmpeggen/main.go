package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/emitter"
	"github.com/bnema/purego-ffmpeg/cmd/ffmpeggen/internal/parser"
)

type config struct {
	headersDir string
	rawDir     string
	publicDir  string
}

func (c config) validate() error {
	if c.headersDir == "" || c.rawDir == "" || c.publicDir == "" {
		return fmt.Errorf("--headers-dir, --raw-dir, and --public-dir are required")
	}
	return nil
}

// headerSpec describes a header file to parse.
type headerSpec struct {
	path  string // relative to headersDir, e.g., "libavutil/frame.h"
	lib   string // library name: "avutil", "avcodec", etc.
	scope *parser.Scope
}

func main() {
	var cfg config
	flag.StringVar(&cfg.headersDir, "headers-dir", "/usr/include", "FFmpeg include root")
	flag.StringVar(&cfg.rawDir, "raw-dir", "ffmpeg/internal/raw", "raw output directory")
	flag.StringVar(&cfg.publicDir, "public-dir", "ffmpeg", "public API output directory")
	flag.Parse()
	if err := run(cfg); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(cfg config) error {
	if err := cfg.validate(); err != nil {
		return err
	}
	if err := os.MkdirAll(cfg.rawDir, 0o755); err != nil {
		return fmt.Errorf("create raw dir: %w", err)
	}
	if err := os.MkdirAll(cfg.publicDir, 0o755); err != nil {
		return fmt.Errorf("create public dir: %w", err)
	}

	specs := v1Scope()
	var registerNames []string

	for _, spec := range specs {
		fullPath := filepath.Join(cfg.headersDir, spec.path)
		header, err := parser.ParseFile(fullPath, spec.lib, spec.scope)
		if err != nil {
			return fmt.Errorf("parse %s: %w", spec.path, err)
		}

		// Derive output file name from header path
		outName := outputName(spec.path)
		header.RegisterName = registerName(spec.lib, spec.path)

		// Skip empty headers
		if len(header.Structs) == 0 && len(header.Functions) == 0 && len(header.Enums) == 0 {
			continue
		}

		// Only register headers that have functions to bind
		if len(header.Functions) > 0 {
			registerNames = append(registerNames, header.RegisterName)
		}

		// Emit raw
		rawCode, err := emitter.EmitRaw(header)
		if err != nil {
			return fmt.Errorf("emit raw %s: %w", outName, err)
		}
		if err := os.WriteFile(filepath.Join(cfg.rawDir, outName), []byte(rawCode), 0o644); err != nil {
			return fmt.Errorf("write raw %s: %w", outName, err)
		}

		// Emit public
		pubData := emitter.BuildPublicFileData(header)
		pubCode, err := emitter.EmitPublic(pubData)
		if err != nil {
			return fmt.Errorf("emit public %s: %w", outName, err)
		}
		if err := os.WriteFile(filepath.Join(cfg.publicDir, outName), []byte(pubCode), 0o644); err != nil {
			return fmt.Errorf("write public %s: %w", outName, err)
		}

		fmt.Printf("  %s → %s (%d structs, %d funcs, %d enums)\n",
			spec.path, outName,
			len(header.Structs), len(header.Functions), len(header.Enums))
	}

	// Generate register.go aggregator
	return writeRegisterAggregator(cfg.rawDir, registerNames)
}

// v1Scope returns the header specs for V1 scope.
func v1Scope() []headerSpec {
	return []headerSpec{
		// libavutil
		{
			path: "libavutil/error.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("av_strerror"),
			},
		},
		{
			path: "libavutil/frame.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("av_frame_alloc", "av_frame_free", "av_frame_unref",
					"av_frame_ref", "av_frame_clone", "av_frame_get_buffer",
					"av_frame_make_writable", "av_frame_is_writable"),
				Enums: setOf("AVFrameSideDataType"),
			},
		},
		{
			path: "libavutil/pixfmt.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Enums: setOf("AVPixelFormat"),
			},
		},
		{
			path: "libavutil/samplefmt.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("av_get_sample_fmt_name", "av_get_bytes_per_sample",
					"av_sample_fmt_is_planar"),
				Enums: setOf("AVSampleFormat"),
			},
		},
		{
			path: "libavutil/avutil.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("avutil_version"),
				Enums:     setOf("AVMediaType"),
			},
		},
		{
			path: "libavutil/rational.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Structs:   setOf("AVRational"),
				Functions: setOf("av_rescale_q"),
			},
		},
		{
			path: "libavutil/mathematics.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("av_rescale_q", "av_rescale_q_rnd"),
			},
		},
		{
			path: "libavutil/dict.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Structs:   setOf("AVDictionaryEntry"),
				Functions: setOf("av_dict_get", "av_dict_set", "av_dict_free", "av_dict_count"),
			},
		},
		{
			path: "libavutil/log.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("av_log_set_level", "av_log_get_level"),
			},
		},
		{
			path: "libavutil/imgutils.h",
			lib:  "avutil",
			scope: &parser.Scope{
				Functions: setOf("av_image_get_buffer_size"),
			},
		},
		// libavcodec
		{
			path: "libavcodec/avcodec.h",
			lib:  "avcodec",
			scope: &parser.Scope{
				Functions: setOf(
					"avcodec_alloc_context3", "avcodec_free_context",
					"avcodec_open2",
					"avcodec_parameters_to_context", "avcodec_parameters_from_context",
					"avcodec_send_packet", "avcodec_receive_frame",
					"avcodec_send_frame", "avcodec_receive_packet",
					"avcodec_flush_buffers",
					"avcodec_is_open",
				),
			},
		},
		{
			path: "libavcodec/codec.h",
			lib:  "avcodec",
			scope: &parser.Scope{
				Functions: setOf("avcodec_find_decoder", "avcodec_find_encoder",
					"avcodec_find_decoder_by_name", "avcodec_find_encoder_by_name"),
			},
		},
		{
			path: "libavcodec/packet.h",
			lib:  "avcodec",
			scope: &parser.Scope{
				Functions: setOf("av_packet_alloc", "av_packet_free", "av_packet_unref",
					"av_packet_ref", "av_packet_rescale_ts"),
			},
		},
		// libavformat
		{
			path: "libavformat/avformat.h",
			lib:  "avformat",
			scope: &parser.Scope{
				Functions: setOf(
					"avformat_version",
					"avformat_alloc_context", "avformat_free_context",
					"avformat_open_input", "avformat_close_input",
					"avformat_find_stream_info",
					"av_find_best_stream",
					"av_read_frame",
					"avformat_alloc_output_context2",
					"avformat_new_stream",
					"avformat_write_header",
					"av_interleaved_write_frame",
					"av_write_trailer",
				),
				Structs: setOf("AVOutputFormat", "AVInputFormat", "AVStream"),
			},
		},
		{
			path: "libavformat/avio.h",
			lib:  "avformat",
			scope: &parser.Scope{
				Functions: setOf("avio_open", "avio_closep"),
			},
		},
		// libswscale
		{
			path: "libswscale/swscale.h",
			lib:  "swscale",
			scope: &parser.Scope{
				Functions: setOf("sws_getContext", "sws_scale", "sws_freeContext",
					"swscale_version"),
			},
		},
		// libswresample
		{
			path: "libswresample/swresample.h",
			lib:  "swresample",
			scope: &parser.Scope{
				Functions: setOf("swr_alloc", "swr_init", "swr_free", "swr_convert",
					"swresample_version"),
			},
		},
	}
}

func setOf(names ...string) map[string]bool {
	m := make(map[string]bool, len(names))
	for _, n := range names {
		m[n] = true
	}
	return m
}

// outputName derives the Go output filename from a header path.
func outputName(headerPath string) string {
	// "libavformat/avformat.h" -> "avformat.go"
	base := filepath.Base(headerPath)
	name := strings.TrimSuffix(base, ".h")
	return name + ".go"
}

// registerName derives a unique Go register function name.
func registerName(lib, headerPath string) string {
	base := filepath.Base(headerPath)
	name := strings.TrimSuffix(base, ".h")
	// PascalCase
	parts := strings.Split(name, "_")
	var sb strings.Builder
	sb.WriteString("Register")
	for _, p := range parts {
		if p == "" {
			continue
		}
		sb.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return sb.String()
}

// writeRegisterAggregator generates register.go for the raw package.
func writeRegisterAggregator(dir string, names []string) error {
	sort.Strings(names)
	// Deduplicate
	unique := make([]string, 0, len(names))
	seen := make(map[string]bool)
	for _, n := range names {
		if !seen[n] {
			seen[n] = true
			unique = append(unique, n)
		}
	}

	var sb strings.Builder
	sb.WriteString("// Code generated by ffmpeggen. DO NOT EDIT.\n\n")
	sb.WriteString("package raw\n\n")
	sb.WriteString("// Handles holds all FFmpeg shared library handles.\n")
	sb.WriteString("type Handles struct {\n")
	sb.WriteString("\tAvutil     uintptr\n")
	sb.WriteString("\tAvcodec    uintptr\n")
	sb.WriteString("\tAvformat   uintptr\n")
	sb.WriteString("\tSwscale    uintptr\n")
	sb.WriteString("\tSwresample uintptr\n")
	sb.WriteString("}\n\n")
	sb.WriteString("// Register loads all FFmpeg C API symbols from the shared libraries.\n")
	sb.WriteString("func Register(h Handles) {\n")

	// Group register calls by library
	for _, name := range unique {
		lib := guessLibFromRegisterName(name)
		sb.WriteString(fmt.Sprintf("\t%s(h.%s)\n", name, lib))
	}
	sb.WriteString("}\n")

	return os.WriteFile(filepath.Join(dir, "register.go"), []byte(sb.String()), 0o644)
}

// guessLibFromRegisterName maps a register function name to a Handles field.
func guessLibFromRegisterName(name string) string {
	lower := strings.ToLower(name)
	switch {
	case strings.Contains(lower, "swscale") || strings.Contains(lower, "sws"):
		return "Swscale"
	case strings.Contains(lower, "swresample") || strings.Contains(lower, "swr"):
		return "Swresample"
	case strings.Contains(lower, "avformat") || strings.Contains(lower, "avio"):
		return "Avformat"
	case strings.Contains(lower, "avcodec") || strings.Contains(lower, "codec") || strings.Contains(lower, "packet"):
		return "Avcodec"
	default:
		return "Avutil"
	}
}
