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
	_dbif_gate__getError  unsafe.Pointer
	_dbif_gate_dbIf_Open  unsafe.Pointer
	_dbif_gate_dbIf_Close unsafe.Pointer
	_dbif_gate_dbBatch_Do unsafe.Pointer
)

func load_dbif(dllPath string) (err error) {
	handle := C.dlopen(C.CString(dllPath), C.RTLD_NOW)
	if uintptr(handle) == 0 {
		return fmt.Errorf("Load %s failed", dllPath)
	}
	_dbif_gate__getError = C.dlsym(handle, C.CString("GetError"))
	if uintptr(_helloif_gate__getError) == 0 {
		return fmt.Errorf("Failed to load function: %s", "GetError")
	}
	_dbif_gate_dbIf_Open = C.dlsym(handle, C.CString("dbIf_Open"))
	if uintptr(_dbif_gate_dbIf_Open) == 0 {
		return fmt.Errorf("Failed to load: %s", "dbIf_Open")
	}
	_dbif_gate_dbIf_Close = C.dlsym(handle, C.CString("dbIf_Close"))
	if uintptr(_dbif_gate_dbIf_Close) == 0 {
		return fmt.Errorf("Failed to load: %s", "dbIf_Close")
	}
	_dbif_gate_dbBatch_Do = C.dlsym(handle, C.CString("dbBatch_Do"))
	if uintptr(_dbif_gate_dbBatch_Do) == 0 {
		return fmt.Errorf("Failed to load: %s", "dbBatch_Do")
	}

	getcrc := C.dlsym(handle, C.CString("GetCRC"))
	if uintptr(getcrc) == 0 {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	C.invokeCrc(getcrc, unsafe.Pointer(&crc))
	if crc != 0x98b5330a8380a2f0 {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x98b5330a8380a2f0", crc)
	}
	return nil
}

func getError_dbif(errPtr unsafe.Pointer) error {
	errText := make([]byte, 0, 512)
	C.invokeGetError(_dbif_gate__getError, unsafe.Pointer(&errText), errPtr)
	return errors.New(string(errText))
}

func (r *dbIf) Open() (err error) {
	errPtr := C.invoke(_dbif_gate_dbIf_Open, unsafe.Pointer(r),
		24)
	if uintptr(errPtr) != 0 {
		return getError_helloif(errPtr)
	}
	return nil
}

func (r *dbIf) Close() (err error) {
	errPtr := C.invoke(_dbif_gate_dbIf_Close, unsafe.Pointer(r),
		24)
	if uintptr(errPtr) != 0 {
		return getError_helloif(errPtr)
	}
	return nil
}

func (r *dbBatch) Do() (err error) {
	errPtr := C.invoke(_dbif_gate_dbBatch_Do, unsafe.Pointer(r),
		24)
	if uintptr(errPtr) != 0 {
		return getError_helloif(errPtr)
	}
	return nil
}
