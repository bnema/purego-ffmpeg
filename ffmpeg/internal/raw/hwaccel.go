package raw

import (
	"unsafe"

	"github.com/ebitengine/purego"
)

// --- Device context ---

// AVHwdeviceCtxCreate wraps av_hwdevice_ctx_create. The first argument is a
// double pointer (AVBufferRef**) — the caller must pass unsafe.Pointer(&buf)
// so that the function can write the newly-created device context back.
var AVHwdeviceCtxCreate func(unsafe.Pointer, int32, *byte, unsafe.Pointer, int32) int32

var AVHwdeviceFindTypeByName func(*byte) int32

var AVHwdeviceGetTypeName func(int32) *byte

var AVHwdeviceIterateTypes func(int32) int32

// --- Frame context ---

var AVHwframeCtxAlloc func(unsafe.Pointer) unsafe.Pointer

var AVHwframeCtxInit func(unsafe.Pointer) int32

var AVHwframeTransferData func(unsafe.Pointer, unsafe.Pointer, int32) int32

// --- Buffer reference counting ---

var AVBufferRef func(unsafe.Pointer) unsafe.Pointer

// AVBufferUnref wraps av_buffer_unref. The argument is a double pointer
// (AVBufferRef**); the C function frees the reference and sets *buf to NULL.
var AVBufferUnref func(unsafe.Pointer)

func RegisterHwaccel(handle uintptr) {
	purego.RegisterLibFunc(&AVHwdeviceCtxCreate, handle, "av_hwdevice_ctx_create")
	purego.RegisterLibFunc(&AVHwdeviceFindTypeByName, handle, "av_hwdevice_find_type_by_name")
	purego.RegisterLibFunc(&AVHwdeviceGetTypeName, handle, "av_hwdevice_get_type_name")
	purego.RegisterLibFunc(&AVHwdeviceIterateTypes, handle, "av_hwdevice_iterate_types")
	purego.RegisterLibFunc(&AVHwframeCtxAlloc, handle, "av_hwframe_ctx_alloc")
	purego.RegisterLibFunc(&AVHwframeCtxInit, handle, "av_hwframe_ctx_init")
	purego.RegisterLibFunc(&AVHwframeTransferData, handle, "av_hwframe_transfer_data")
	purego.RegisterLibFunc(&AVBufferRef, handle, "av_buffer_ref")
	purego.RegisterLibFunc(&AVBufferUnref, handle, "av_buffer_unref")
}
