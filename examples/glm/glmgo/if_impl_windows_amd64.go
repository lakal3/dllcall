//
package main

// Generated file. Not not edit

import "syscall"
import "unsafe"
import "errors"
import "fmt"

var (
	_if_gate__getError                    uintptr
	_if_gate_MultiplyVectors_Multiply     uintptr
	_if_gate_MultiplyVectors_FastMultiply uintptr
)

func _if_fastcall(trap uintptr, ptr uintptr, size uintptr) (errPtr uintptr)

func load_if(dllPath string) (err error) {
	if _if_gate__getError != 0 {
		return nil
	}
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_if_gate_MultiplyVectors_Multiply, err = syscall.GetProcAddress(dll, "MultiplyVectors_Multiply")
	if err != nil {
		return fmt.Errorf("%s: %v", "MultiplyVectors_Multiply", err)
	}
	_if_gate_MultiplyVectors_FastMultiply, err = syscall.GetProcAddress(dll, "MultiplyVectors_FastMultiply")
	if err != nil {
		return fmt.Errorf("%s: %v", "MultiplyVectors_FastMultiply", err)
	}

	getcrc, err := syscall.GetProcAddress(dll, "GetCRC")
	if err != nil {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	syscall.Syscall(getcrc, 1, uintptr(unsafe.Pointer(&crc)), 0, 0)
	if crc != 0xfad9b164355a2f76 {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0xfad9b164355a2f76", crc)
	}
	_if_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
	if err != nil {
		return err
	}
	return nil
}

func if_getError(rc uintptr) error {
	errText := make([]byte, 0, 512)
	syscall.Syscall(_if_gate__getError, 2, rc, uintptr(unsafe.Pointer(&errText)), 0)
	return errors.New(string(errText))
}

func (r *MultiplyVectors) Multiply() (err error) {
	rc, _, _ := syscall.Syscall(_if_gate_MultiplyVectors_Multiply, 2, uintptr(unsafe.Pointer(r)),
		uintptr(152), 0)
	if rc != 0 {
		return if_getError(rc)
	}
	return nil
}

func (r *MultiplyVectors) FastMultiply() (err error) {
	rc := _if_fastcall(_if_gate_MultiplyVectors_FastMultiply, uintptr(unsafe.Pointer(r)),
		uintptr(152))
	if rc != 0 {
		return if_getError(rc)
	}
	return nil
}
