package capi

import "github.com/bnema/purego"

// tryRegisterLibFunc registers a C function symbol if it exists in the library.
// Unlike purego.RegisterLibFunc, it silently skips symbols not found in the
// shared library, which is essential because optional FFmpeg APIs may vary
// across builds and configurations.
func tryRegisterLibFunc(fptr any, handle uintptr, name string) {
	if _, err := purego.Dlsym(handle, name); err != nil {
		return // symbol not available in this build
	}
	purego.RegisterLibFunc(fptr, handle, name)
}
