package syscall

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>

void *invoke2(void *addr, void *arg1, void *arg2) {
void *(*ptr)(void *arg, void * size) = addr;
	return (*ptr)(arg1, arg2);
}

void *invoke1(void *addr, void *arg1)  {
	unsigned long long (*ptr)(void *arg1) = addr;
	(*ptr)(arg1);
}
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type Handle uintptr

func LoadLibrary(libraryPath string) (Handle, error) {
	handle := C.dlopen(C.CString(libraryPath), C.RTLD_NOW)
	if uintptr(handle) == 0 {
		return 0, fmt.Errorf("Load %s failed", libraryPath)
	}
	return Handle(handle), nil
}

func GetProcAddress(dll Handle, methodName string) (trap uintptr, err error) {
	mh := C.dlsym(unsafe.Pointer(dll), C.CString(methodName))
	if uintptr(mh) == 0 {
		return 0, fmt.Errorf("Failed to load function: %s", methodName)
	}
	return uintptr(mh), nil
}

func Syscall(trap uintptr, nargs int, a1 uintptr, a2 uintptr, a3 uintptr) uintptr {
	switch nargs {
	case 1:
		return uintptr(C.invoke1(unsafe.Pointer(trap), unsafe.Pointer(a1)))
	case 2:
		return uintptr(C.invoke2(unsafe.Pointer(trap), unsafe.Pointer(a1), unsafe.Pointer(a2)))
	}
	panic("Nargs must be 1 or 2")
}
