// internal/loader/loader.go
package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/bnema/purego"
)

// Handles holds all FFmpeg shared library handles.
type Handles struct {
	Avutil     uintptr
	Avcodec    uintptr
	Avformat   uintptr
	Swscale    uintptr
	Swresample uintptr
	Avfilter   uintptr
}

// SOVersions specifies the SO major version numbers for each library.
type SOVersions struct {
	Avutil     string
	Avcodec    string
	Avformat   string
	Swscale    string
	Swresample string
	Avfilter   string
}

// Option configures the Load call.
type Option func(*config)

type config struct {
	dir        string
	soVersions SOVersions
}

// defaultSOVersions corresponds to FFmpeg 7.x releases.
// Override with WithSOVersions for other FFmpeg versions.
var defaultSOVersions = SOVersions{
	Avutil:     "60",
	Avcodec:    "62",
	Avformat:   "62",
	Swscale:    "9",
	Swresample: "6",
	Avfilter:   "11",
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
		{"libavfilter", cfg.soVersions.Avfilter},
	}

	var handles [6]uintptr
	for i, l := range libs {
		soName := libFileName(l.name, l.soVer)
		var fullPath string
		if dir != "" {
			fullPath = filepath.Join(dir, soName)
		} else {
			fullPath = soName
		}
		handle, err := purego.Dlopen(fullPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			// Close previously opened handles to avoid resource leaks.
			for j := 0; j < i; j++ {
				purego.Dlclose(handles[j]) //nolint:errcheck
			}
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
		Avfilter:   handles[5],
	}, nil
}

// libFileName returns the platform-specific shared library filename for the
// given base name and major version number.
func libFileName(name, version string) string {
	switch runtime.GOOS {
	case "darwin":
		return name + "." + version + ".dylib"
	case "windows":
		return name + "-" + version + ".dll"
	default:
		return name + ".so." + version
	}
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
