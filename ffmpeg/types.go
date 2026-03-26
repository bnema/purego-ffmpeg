package ffmpeg

import portin "github.com/bnema/purego-ffmpeg/internal/ports/in"

// AVRational aliases portin.AVRational so that the public API uses the same
// type as the port interfaces without re-defining it.
type AVRational = portin.AVRational
