package mocks

import portin "github.com/bnema/purego-ffmpeg/internal/ports/in"

// AVRational aliases portin.AVRational so that mock methods satisfy
// the public interfaces when passing this value type.
type AVRational = portin.AVRational
