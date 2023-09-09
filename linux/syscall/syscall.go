package syscall

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>

void *invoke3(void *addr, void *arg1, void *arg2, void *arg3) {
void *(*ptr)(void *arg1, void * arg2, void *arg3) = addr;
	return (*ptr)(arg1, arg2, arg3);
}

void *invoke2(void *addr, void *arg1, void *arg2) {
void *(*ptr)(void *arg1, void * arg2) = addr;
	return (*ptr)(arg1, arg2);
}

void *invoke1(void *addr, void *arg1)  {
	unsigned long long (*ptr)(void *arg1) = addr;
	(*ptr)(arg1);
}

void *invoke0(void *addr)  {
	unsigned long long (*ptr)() = addr;
	(*ptr)();
}
*/
import "C"

import (
	"fmt"
	"strings"
	"unsafe"
)

type Handle uintptr

// LoadLibrary load shared object file (.so) and gives a handle to it.
// If libraryPath name looks like Windows DLL name (is has .dll extension) it will be converted to Linux equivalent
// Like someapi.dll -> libsomeapi.so
func LoadLibrary(libraryPath string) (Handle, error) {
	if strings.HasSuffix(strings.ToLower(libraryPath), ".dll") {
		// Convert Windows dll name to Linux equivalent
		libraryPath = "lib" + libraryPath[:len(libraryPath)-4] + ".so"
	}
	handle := C.dlopen(C.CString(libraryPath), C.RTLD_NOW)
	if uintptr(handle) == 0 {
		errText := string(C.GoString(C.dlerror()))
		return 0, fmt.Errorf("Load %s failed: %s", libraryPath, errText)
	}
	return Handle(handle), nil
}

// FreeLibrary releases previously loaded library
func FreeLibrary(dll Handle) {
	C.dlclose(unsafe.Pointer(dll))
}

// Retrieve address to exported function
func GetProcAddress(dll Handle, exportName string) (trap uintptr, err error) {
	mh := C.dlsym(unsafe.Pointer(dll), C.CString(exportName))
	if uintptr(mh) == 0 {
		return 0, fmt.Errorf("Failed to load function: %s", exportName)
	}
	return uintptr(mh), nil
}

// SyscallL invoke function retrieved with GetProcAddress
func SyscallL(trap uintptr, args ...uintptr) uintptr {
	switch len(args) {
	case 0:
		return uintptr(C.invoke0(unsafe.Pointer(trap)))
	case 1:
		return uintptr(C.invoke1(unsafe.Pointer(trap), unsafe.Pointer(args[0])))
	case 2:
		return uintptr(C.invoke2(unsafe.Pointer(trap), unsafe.Pointer(args[0]), unsafe.Pointer(args[1])))
	case 3:
		return uintptr(C.invoke3(unsafe.Pointer(trap), unsafe.Pointer(args[0]), unsafe.Pointer(args[1]), unsafe.Pointer(args[2])))
	}
	panic("len args must be from 0 to 3")
}

// SyscallN with matching signature to Windows one. r1 will always be 0 and err nil
func SyscallN(trap uintptr, args ...uintptr) (r0 uintptr, r1 uintptr, err error) {
	return SyscallL(trap, args...), 0, nil
}
