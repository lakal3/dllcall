// +build windows

//
package main

// Generated file. Not not edit

import "syscall"
import "unsafe"
import "errors"
import "fmt"

var (
	_if_gate__getError                uintptr
	_if_gate_MultiplyVectors_Multiply uintptr
)

func load_if(dllPath string) (err error) {
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_if_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
	if err != nil {
		return err
	}
	_if_gate_MultiplyVectors_Multiply, err = syscall.GetProcAddress(dll, "MultiplyVectors_Multiply")
	if err != nil {
		return fmt.Errorf("%s: %v", "MultiplyVectors_Multiply", err)
	}

	getcrc, err := syscall.GetProcAddress(dll, "GetCRC")
	if err != nil {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	syscall.Syscall(getcrc, 1, uintptr(unsafe.Pointer(&crc)), 0, 0)
	if crc != 0x00698cefaddc526c {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x00698cefaddc526c", crc)
	}
	return nil
}

func getError_if(rc uintptr) error {
	errText := make([]byte, 0, 512)
	syscall.Syscall(_if_gate__getError, 2, rc, uintptr(unsafe.Pointer(&errText)), 0)
	return errors.New(string(errText))
}

func (r *MultiplyVectors) Multiply() (err error) {
	rc, _, _ := syscall.Syscall(_if_gate_MultiplyVectors_Multiply, 2, uintptr(unsafe.Pointer(r)),
		uintptr(152), 0)
	if rc != 0 {
		return getError_if(rc)
	}
	return nil
}
