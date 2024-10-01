package syscall

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>


typedef struct dummyArg {
	long long tmp;
} dummyArg;

void invokeGetError(void *addr, void *errPtr, dummyArg *errResult)  {
	unsigned long long (*ptr)(void *errPtr, dummyArg *errResult) = addr;
	(*ptr)(errPtr, errResult);
}

void invokeCRC(void *addr, unsigned long long *crc)  {
	void (*ptr)(unsigned long long *crc) = addr;
	(*ptr)(crc);
}


void *invokeT(void *addr, dummyArg *arg, size_t argLen) {
	void *(*ptr)(void *arg, size_t argLen) = addr;
	return (*ptr)(arg, argLen);
}

*/
import "C"

import (
	"fmt"
	"runtime"
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

func GetCRC(trap uintptr, crc *uint64) {
	C.invokeCRC(unsafe.Pointer(trap), (*C.ulonglong)(crc))
}

func GetError(trap uintptr, errPtr uintptr, errText *[]byte) {
	C.invokeGetError(unsafe.Pointer(trap), unsafe.Pointer(errPtr), (*C.dummyArg)(unsafe.Pointer(errText)))
}

func SyscallT[T any](trap uintptr, arg *T) (r0 uintptr) {
	r0 = uintptr(C.invokeT(unsafe.Pointer(trap), (*C.dummyArg)(unsafe.Pointer(arg)), C.size_t(unsafe.Sizeof(*arg))))
	runtime.KeepAlive(arg)
	return r0
}
