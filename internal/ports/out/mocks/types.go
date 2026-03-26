package mocks

import portout "github.com/bnema/purego-ffmpeg/internal/ports/out"

// AVRational aliases portout.AVRational so that mock methods satisfy
// the portout.* port interfaces when passing this value type.
type AVRational = portout.AVRational
