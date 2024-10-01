package main

// Generated file. Do not edit

import "github.com/lakal3/dllcall/linux/syscall"
import "errors"
import "fmt"
import "unsafe"

var (
	_if_gate__getError                    uintptr
	_if_gate_MultiplyVectors_Multiply     uintptr
	_if_gate_MultiplyVectors_FastMultiply uintptr
)

func _if_fastcall(trap uintptr, ptr uintptr, size uintptr) (errPtr uintptr)
func _if_fc_alloc() (ret uintptr)

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
	syscall.GetCRC(getcrc, &crc)
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
	errText := make([]byte, 0, 4096)
	syscall.GetError(_if_gate__getError, rc, &errText)
	return errors.New(string(errText))
}

func (r *MultiplyVectors) Multiply() (err error) {
	rc := syscall.SyscallT(_if_gate_MultiplyVectors_Multiply, r)
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
