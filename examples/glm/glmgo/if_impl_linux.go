// +build linux

//
package main

// Generated file. Not not edit

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>

void *invoke(void *addr, void *arg, long long size) {
void *(*ptr)(void *arg, long long size) = addr;
	return (*ptr)(arg, size);
}
void invokeGetError(void *addr, void *errPtr, void *msgSlice)  {
void (*ptr)(void * errPtr, void *msgSlice) = addr;
	(*ptr)(errPtr, msgSlice);
}

void invokeCrc(void *addr, void *crc)  {
	unsigned long long (*ptr)(void *crc) = addr;
	(*ptr)(crc);
}
*/
import "C"
import "unsafe"
import "errors"
import "fmt"

var (
	_if_gate__getError                unsafe.Pointer
	_if_gate_MultiplyVectors_Multiply unsafe.Pointer
)

func load_if(dllPath string) (err error) {
	handle := C.dlopen(C.CString(dllPath), C.RTLD_NOW)
	if uintptr(handle) == 0 {
		return fmt.Errorf("Load %s failed", dllPath)
	}
	_if_gate__getError = C.dlsym(handle, C.CString("GetError"))
	if uintptr(_helloif_gate__getError) == 0 {
		return fmt.Errorf("Failed to load function: %s", "GetError")
	}
	_if_gate_MultiplyVectors_Multiply = C.dlsym(handle, C.CString("MultiplyVectors_Multiply"))
	if uintptr(_if_gate_MultiplyVectors_Multiply) == 0 {
		return fmt.Errorf("Failed to load: %s", "MultiplyVectors_Multiply")
	}

	getcrc := C.dlsym(handle, C.CString("GetCRC"))
	if uintptr(getcrc) == 0 {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	C.invokeCrc(getcrc, unsafe.Pointer(&crc))
	if crc != 0x00698cefaddc526c {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x00698cefaddc526c", crc)
	}
	return nil
}

func getError_if(errPtr unsafe.Pointer) error {
	errText := make([]byte, 0, 512)
	C.invokeGetError(_if_gate__getError, unsafe.Pointer(&errText), errPtr)
	return errors.New(string(errText))
}

func (r *MultiplyVectors) Multiply() (err error) {
	errPtr := C.invoke(_if_gate_MultiplyVectors_Multiply, unsafe.Pointer(r),
		152)
	if uintptr(errPtr) != 0 {
		return getError_helloif(errPtr)
	}
	return nil
}
