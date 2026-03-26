// internal/core/strings.go
package core

import (
	"sync"
	"unsafe"
)

var pinnedStrings sync.Map // *byte → []byte — prevents GC collection

// CString converts a Go string to a null-terminated C string.
// The caller must call FreeCString when done to allow GC collection.
func CString(s string) *byte {
	b := make([]byte, len(s)+1)
	copy(b, s)
	ptr := &b[0]
	pinnedStrings.Store(ptr, b) // prevent GC from collecting b
	return ptr
}

// GoString converts a null-terminated C string to a Go string.
// The input must be a valid null-terminated byte sequence.
// Passing a non-null-terminated pointer will cause an unbounded read.
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

// FreeCString releases the reference to the backing array,
// allowing the GC to collect it.
func FreeCString(p *byte) {
	if p != nil {
		pinnedStrings.Delete(p)
	}
}
