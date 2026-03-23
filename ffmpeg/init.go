package ffmpeg

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/ebitengine/purego"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// Default SO version numbers for FFmpeg 7.x / 8.x.
var defaultSOVersions = SOVersions{
	Avutil:     "60",
	Avcodec:    "62",
	Avformat:   "62",
	Swscale:    "9",
	Swresample: "6",
}

// SOVersions specifies the SO major version numbers for each library.
type SOVersions struct {
	Avutil     string
	Avcodec    string
	Avformat   string
	Swscale    string
	Swresample string
}

// Option configures the Init call.
type Option func(*initConfig)

type initConfig struct {
	dir        string
	soVersions SOVersions
}

// WithDir sets the directory to search for FFmpeg shared libraries.
func WithDir(dir string) Option {
	return func(c *initConfig) { c.dir = dir }
}

// WithSOVersions overrides the default SO version numbers.
func WithSOVersions(v SOVersions) Option {
	return func(c *initConfig) { c.soVersions = v }
}

var initOnce sync.Once

// Init loads all FFmpeg shared libraries and registers symbols.
// Must be called before using any FFmpeg functions.
func Init(opts ...Option) error {
	cfg := initConfig{soVersions: defaultSOVersions}
	for _, o := range opts {
		o(&cfg)
	}

	var initErr error
	initOnce.Do(func() {
		initErr = doInit(cfg)
	})
	return initErr
}

func doInit(cfg initConfig) error {
	dir := resolveDir(cfg.dir)

	// Load libraries in dependency order.
	libs := []struct {
		name    string
		soVer   string
		handleP *uintptr
	}{
		{"libavutil", cfg.soVersions.Avutil, new(uintptr)},
		{"libswresample", cfg.soVersions.Swresample, new(uintptr)},
		{"libavcodec", cfg.soVersions.Avcodec, new(uintptr)},
		{"libavformat", cfg.soVersions.Avformat, new(uintptr)},
		{"libswscale", cfg.soVersions.Swscale, new(uintptr)},
	}

	for _, lib := range libs {
		soName := lib.name + ".so." + lib.soVer
		var fullPath string
		if dir != "" {
			fullPath = filepath.Join(dir, soName)
		} else {
			fullPath = soName
		}
		handle, err := purego.Dlopen(fullPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return fmt.Errorf("load %s: %w", soName, err)
		}
		*lib.handleP = handle
	}

	raw.Register(raw.Handles{
		Avutil:     *libs[0].handleP,
		Swresample: *libs[1].handleP,
		Avcodec:    *libs[2].handleP,
		Avformat:   *libs[3].handleP,
		Swscale:    *libs[4].handleP,
	})

	return nil
}

// resolveDir finds the FFmpeg library directory.
func resolveDir(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if dir := os.Getenv("FFMPEG_DIR"); dir != "" {
		return dir
	}
	// Common system paths — Dlopen with just the soname will search ld.so.conf.
	return ""
}

// Shutdown performs any global cleanup. Currently a no-op but reserved
// for future use (e.g., avformat_network_deinit).
func Shutdown() {}
