# purego-ffmpeg

Go bindings for the FFmpeg C API using [purego](https://github.com/ebitengine/purego) — no cgo required.

> **Early stage** — API is not stable. Currently targets FFmpeg 7.x/8.x on Linux.

## Architecture

Same approach as [purego-cef](https://github.com/bnema/purego-cef): a code generator parses FFmpeg C headers and produces two layers of Go code.

```
purego-ffmpeg/
├── cmd/ffmpeggen/          # Code generator (regex-based header parser + templates)
│   └── internal/
│       ├── parser/         # FFmpeg C header parser
│       ├── model/          # AST models
│       └── emitter/        # Go code generator with templates
├── ffmpeg/                 # Public API (generated + hand-written)
│   ├── internal/raw/       # Generated: C struct layouts + purego symbol registration
│   ├── init.go             # Hand-written: library loading, version validation
│   ├── support.go          # Hand-written: error handling, string helpers
│   └── *.go               # Generated: wrappers, enums, function bindings
└── integration/            # Integration tests (requires FFmpeg runtime)
```

## V1 Scope

Libraries bound:
- **libavutil** — error codes, frame, pixel/sample formats, dictionary, logging, rationals
- **libavcodec** — codec context, encoder/decoder lookup, packet send/receive
- **libavformat** — format context, demuxing, muxing, stream info
- **libswscale** — video scaling/conversion
- **libswresample** — audio resampling

## Requirements

- Go 1.26+
- FFmpeg 7.x or 8.x shared libraries (`libavutil.so.60`, `libavcodec.so.62`, etc.)
- Linux (for now)

## Usage

```go
package main

import (
    "fmt"
    "github.com/bnema/purego-ffmpeg/ffmpeg"
)

func main() {
    if err := ffmpeg.Init(); err != nil {
        panic(err)
    }

    ver := ffmpeg.FormatVersion()
    fmt.Printf("avformat version: %d.%d.%d\n", ver>>16, (ver>>8)&0xFF, ver&0xFF)
}
```

## Building

```bash
CGO_ENABLED=0 go build ./...
CGO_ENABLED=0 go test ./ffmpeg/
```

## Integration Tests

Requires FFmpeg shared libraries at runtime:

```bash
go test -tags=ffmpeg_integration ./integration/ -v
```

## Regenerating Bindings

```bash
go run ./cmd/ffmpeggen/ \
  --headers-dir=/usr/include \
  --raw-dir=ffmpeg/internal/raw \
  --public-dir=ffmpeg
```

## License

MIT
