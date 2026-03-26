package capi

import portout "github.com/bnema/purego-ffmpeg/internal/ports/out"

// AVRational aliases portout.AVRational so that CAPI adapter methods
// satisfy the portout.* port interfaces when passing this value type.
type AVRational = portout.AVRational
