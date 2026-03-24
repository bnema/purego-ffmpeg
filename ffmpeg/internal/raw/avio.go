package raw

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

var AvioOpen func(unsafe.Pointer, *byte, int32) int32

var AvioClosep func(unsafe.Pointer) int32

// avio_alloc_context(buffer, buffer_size, write_flag, opaque, read_packet, write_packet, seek)
// Callback parameters are uintptr (from purego.NewCallback).
var AvioAllocContext func(unsafe.Pointer, int32, int32, unsafe.Pointer, uintptr, uintptr, uintptr) unsafe.Pointer

var AvioContextFree func(unsafe.Pointer)

func RegisterAvio(handle uintptr) {
	purego.RegisterLibFunc(&AvioOpen, handle, "avio_open")
	purego.RegisterLibFunc(&AvioClosep, handle, "avio_closep")
	purego.RegisterLibFunc(&AvioAllocContext, handle, "avio_alloc_context")
	purego.RegisterLibFunc(&AvioContextFree, handle, "avio_context_free")
}
