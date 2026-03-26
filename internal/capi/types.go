package capi

import out "github.com/bnema/purego-ffmpeg/internal/ports/out"

// AVRational aliases out.AVRational so that CAPI adapter methods
// satisfy the out.* port interfaces when passing this value type.
type AVRational = out.AVRational
