package in

import portout "github.com/bnema/purego-ffmpeg/internal/ports/out"

// AVRational aliases portout.AVRational so that port-in interfaces use the
// same concrete type as the outbound port layer. This is necessary
// because AVRational appears in function signatures across packages
// and Go requires type identity for interface satisfaction.
type AVRational = portout.AVRational
