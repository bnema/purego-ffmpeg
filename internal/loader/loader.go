// internal/loader/loader.go
package loader

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ebitengine/purego"
)

// Handles holds all FFmpeg shared library handles.
type Handles struct {
	Avutil     uintptr
	Avcodec    uintptr
	Avformat   uintptr
	Swscale    uintptr
	Swresample uintptr
}

// SOVersions specifies the SO major version numbers for each library.
type SOVersions struct {
	Avutil     string
	Avcodec    string
	Avformat   string
	Swscale    string
	Swresample string
}

// Option configures the Load call.
type Option func(*config)

type config struct {
	dir        string
	soVersions SOVersions
}

var defaultSOVersions = SOVersions{
	Avutil:     "60",
	Avcodec:    "62",
	Avformat:   "62",
	Swscale:    "9",
	Swresample: "6",
}

// WithDir sets the directory to search for FFmpeg shared libraries.
func WithDir(dir string) Option {
	return func(c *config) { c.dir = dir }
}

// WithSOVersions overrides the default SO version numbers.
func WithSOVersions(v SOVersions) Option {
	return func(c *config) { c.soVersions = v }
}

// Load opens all FFmpeg shared libraries and returns their handles.
// Libraries are loaded in dependency order.
func Load(opts ...Option) (Handles, error) {
	cfg := config{soVersions: defaultSOVersions}
	for _, o := range opts {
		o(&cfg)
	}
	dir := resolveDir(cfg.dir)

	type lib struct {
		name  string
		soVer string
	}
	libs := []lib{
		{"libavutil", cfg.soVersions.Avutil},
		{"libswresample", cfg.soVersions.Swresample},
		{"libavcodec", cfg.soVersions.Avcodec},
		{"libavformat", cfg.soVersions.Avformat},
		{"libswscale", cfg.soVersions.Swscale},
	}

	var handles [5]uintptr
	for i, l := range libs {
		soName := l.name + ".so." + l.soVer
		var fullPath string
		if dir != "" {
			fullPath = filepath.Join(dir, soName)
		} else {
			fullPath = soName
		}
		handle, err := purego.Dlopen(fullPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return Handles{}, fmt.Errorf("load %s: %w", soName, err)
		}
		handles[i] = handle
	}

	return Handles{
		Avutil:     handles[0],
		Swresample: handles[1],
		Avcodec:    handles[2],
		Avformat:   handles[3],
		Swscale:    handles[4],
	}, nil
}

func resolveDir(explicit string) string {
	if explicit != "" {
		return explicit
	}
	if dir := os.Getenv("FFMPEG_DIR"); dir != "" {
		return dir
	}
	return ""
}
