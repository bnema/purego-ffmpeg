package ffmpeg

import (
	"runtime"
	"unsafe"

	"github.com/bnema/purego-ffmpeg/ffmpeg/internal/raw"
)

// Hardware device types from enum AVHWDeviceType.
// Values sourced from libavutil/hwcontext.h.
const (
	AV_HWDEVICE_TYPE_NONE         int32 = 0
	AV_HWDEVICE_TYPE_VDPAU        int32 = 1
	AV_HWDEVICE_TYPE_CUDA         int32 = 2
	AV_HWDEVICE_TYPE_VAAPI        int32 = 3
	AV_HWDEVICE_TYPE_DXVA2        int32 = 4
	AV_HWDEVICE_TYPE_QSV          int32 = 5
	AV_HWDEVICE_TYPE_VIDEOTOOLBOX int32 = 6
	AV_HWDEVICE_TYPE_D3D11VA      int32 = 7
	AV_HWDEVICE_TYPE_DRM          int32 = 8
	AV_HWDEVICE_TYPE_OPENCL       int32 = 9
	AV_HWDEVICE_TYPE_MEDIACODEC   int32 = 10
	AV_HWDEVICE_TYPE_VULKAN       int32 = 11
	AV_HWDEVICE_TYPE_D3D12VA      int32 = 12
	AV_HWDEVICE_TYPE_AMF          int32 = 13
	AV_HWDEVICE_TYPE_OHCODEC      int32 = 14
)

// HWDeviceCtxCreate creates a hardware device context.
// deviceType: AV_HWDEVICE_TYPE_VAAPI, AV_HWDEVICE_TYPE_CUDA, etc.
// device: device path (e.g., "/dev/dri/renderD128") or empty for auto-detect.
// Returns the device context buffer ref, or error.
func HWDeviceCtxCreate(deviceType int32, device string) (unsafe.Pointer, error) {
	var deviceCtx unsafe.Pointer
	var deviceC *byte
	var deviceBuf []byte
	if device != "" {
		deviceC, deviceBuf = cString(device)
		defer runtime.KeepAlive(deviceBuf)
	}
	ret := raw.AVHwdeviceCtxCreate(unsafe.Pointer(&deviceCtx), deviceType, deviceC, nil, 0)
	if err := avErr(ret); err != nil {
		return nil, err
	}
	return deviceCtx, nil
}

// HWDeviceFindTypeByName looks up a device type by name ("vaapi", "cuda", "vulkan").
// Returns AV_HWDEVICE_TYPE_NONE if not found.
func HWDeviceFindTypeByName(name string) int32 {
	nameC, nameBuf := cString(name)
	defer runtime.KeepAlive(nameBuf)
	return raw.AVHwdeviceFindTypeByName(nameC)
}

// HWDeviceGetTypeName returns the name for a device type constant.
// Returns empty string if the type is not valid.
func HWDeviceGetTypeName(deviceType int32) string {
	return goString(raw.AVHwdeviceGetTypeName(deviceType))
}

// HWDeviceIterateTypes iterates available device types.
// Pass AV_HWDEVICE_TYPE_NONE to start. Returns AV_HWDEVICE_TYPE_NONE when done.
func HWDeviceIterateTypes(prev int32) int32 {
	return raw.AVHwdeviceIterateTypes(prev)
}

// HWFrameCtxAlloc allocates a hardware frame context from a device context.
// Returns nil on failure.
func HWFrameCtxAlloc(deviceCtx unsafe.Pointer) unsafe.Pointer {
	return raw.AVHwframeCtxAlloc(deviceCtx)
}

// HWFrameCtxInit initializes a hardware frame context after its fields have
// been set. Returns 0 on success, negative AVERROR on failure.
func HWFrameCtxInit(ref unsafe.Pointer) int32 {
	return raw.AVHwframeCtxInit(ref)
}

// HWFrameTransferData transfers frame data between hardware and software
// surfaces. At least one of dst/src must have an AVHWFramesContext attached.
// Returns 0 on success, negative AVERROR on failure.
func HWFrameTransferData(dst, src unsafe.Pointer, flags int32) int32 {
	return raw.AVHwframeTransferData(dst, src, flags)
}

// BufferRef creates a new reference to a buffer.
// Returns nil on failure.
func BufferRef(buf unsafe.Pointer) unsafe.Pointer {
	return raw.AVBufferRef(buf)
}

// BufferUnref frees a buffer reference and sets the pointer to nil.
func BufferUnref(buf unsafe.Pointer) {
	raw.AVBufferUnref(buf)
}
