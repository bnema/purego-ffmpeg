package raw

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

var AVHwframeGetBuffer func(unsafe.Pointer, unsafe.Pointer, int32) int32

func RegisterFrameExtra(handle uintptr) {
	purego.RegisterLibFunc(&AVHwframeGetBuffer, handle, "av_hwframe_get_buffer")
}
