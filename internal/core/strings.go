// internal/core/strings.go
package core

import (
	"unsafe"
)

// CString converts a Go string to a null-terminated C string.
// The caller must call FreeCString when done.
func CString(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

// GoString converts a null-terminated C string to a Go string.
func GoString(p *byte) string {
	if p == nil {
		return ""
	}
	ptr := unsafe.Pointer(p)
	var length int
	for {
		if *(*byte)(unsafe.Add(ptr, length)) == 0 {
			break
		}
		length++
	}
	return string(unsafe.Slice(p, length))
}

// FreeCString is a no-op — the Go GC handles the memory.
// Provided for API symmetry and future compatibility.
func FreeCString(_ *byte) {}
