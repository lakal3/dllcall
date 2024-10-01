package main

// Generated file. Do not edit

import "syscall"
import "unsafe"
import "errors"
import "fmt"

var (
	_fibon_if_gate__getError                  uintptr
	_fibon_if_gate_calcFibonacci_calc         uintptr
	_fibon_if_gate_fastcalcFibonacci_fastCalc uintptr
)

func _fibon_if_fastcall(trap uintptr, ptr uintptr, size uintptr) (errPtr uintptr)

func load_fibon_if(dllPath string) (err error) {
	if _fibon_if_gate__getError != 0 {
		return nil
	}
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_fibon_if_gate_calcFibonacci_calc, err = syscall.GetProcAddress(dll, "calcFibonacci_calc")
	if err != nil {
		return fmt.Errorf("%s: %v", "calcFibonacci_calc", err)
	}
	_fibon_if_gate_fastcalcFibonacci_fastCalc, err = syscall.GetProcAddress(dll, "fastcalcFibonacci_fastCalc")
	if err != nil {
		return fmt.Errorf("%s: %v", "fastcalcFibonacci_fastCalc", err)
	}

	getcrc, err := syscall.GetProcAddress(dll, "GetCRC")
	if err != nil {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	_, _, _ = syscall.SyscallN(getcrc, uintptr(unsafe.Pointer(&crc)))
	if crc != 0x0ec9cea884595c43 {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x0ec9cea884595c43", crc)
	}
	_fibon_if_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
	if err != nil {
		return err
	}
	return nil
}

func fibon_if_getError(rc uintptr) (err error) {
	errText := make([]byte, 0, 4096)
	_, _, _ = syscall.SyscallN(_fibon_if_gate__getError, rc, uintptr(unsafe.Pointer(&errText)))
	return errors.New(string(errText))
}

func (r *calcFibonacci) calc() (err error) {
	rc, _, _ := syscall.SyscallN(_fibon_if_gate_calcFibonacci_calc, uintptr(unsafe.Pointer(r)), uintptr(16))
	if rc != 0 {
		return fibon_if_getError(rc)
	}
	return nil
}

func (r *fastcalcFibonacci) fastCalc() (err error) {
	rc := _fibon_if_fastcall(_fibon_if_gate_fastcalcFibonacci_fastCalc, uintptr(unsafe.Pointer(r)),
		uintptr(16))
	if rc != 0 {
		return fibon_if_getError(rc)
	}
	return nil
}
