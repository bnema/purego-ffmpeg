package raw

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

type MediaType int32

var AvutilVersion func() uint32

var AVMalloc func(uintptr) unsafe.Pointer

var AVFree func(unsafe.Pointer)

var AVOptSet func(unsafe.Pointer, *byte, *byte, int32) int32

var AVOptSetInt func(unsafe.Pointer, *byte, int64, int32) int32

func RegisterAvutil(handle uintptr) {
	purego.RegisterLibFunc(&AvutilVersion, handle, "avutil_version")
	purego.RegisterLibFunc(&AVMalloc, handle, "av_malloc")
	purego.RegisterLibFunc(&AVFree, handle, "av_free")
	purego.RegisterLibFunc(&AVOptSet, handle, "av_opt_set")
	purego.RegisterLibFunc(&AVOptSetInt, handle, "av_opt_set_int")
}
