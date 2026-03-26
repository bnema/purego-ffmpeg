package ffmpeg

import in "github.com/bnema/purego-ffmpeg/internal/ports/in"

// AVRational aliases in.AVRational so that the public API uses the same
// type as the port interfaces without re-defining it.
type AVRational = in.AVRational
