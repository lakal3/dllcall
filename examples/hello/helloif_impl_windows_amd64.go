package main

// Generated file. Not not edit

import "syscall"
import "unsafe"
import "errors"
import "fmt"

var (
	_helloif_gate__getError      uintptr
	_helloif_gate_greeting_Greet uintptr
)

func load_helloif(dllPath string) (err error) {
	if _helloif_gate__getError != 0 {
		return nil
	}
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_helloif_gate_greeting_Greet, err = syscall.GetProcAddress(dll, "greeting_Greet")
	if err != nil {
		return fmt.Errorf("%s: %v", "greeting_Greet", err)
	}

	getcrc, err := syscall.GetProcAddress(dll, "GetCRC")
	if err != nil {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	syscall.SyscallN(getcrc, uintptr(unsafe.Pointer(&crc)))
	if crc != 0x9cc3656ee6911505 {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x9cc3656ee6911505", crc)
	}
	_helloif_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
	if err != nil {
		return err
	}
	return nil
}

func helloif_getError(rc uintptr) error {
	errText := make([]byte, 0, 512)
	syscall.SyscallN(_helloif_gate__getError, rc, uintptr(unsafe.Pointer(&errText)))
	return errors.New(string(errText))
}

func (r *greeting) Greet() (err error) {
	rc, _, _ := syscall.SyscallN(_helloif_gate_greeting_Greet, uintptr(unsafe.Pointer(r)), uintptr(48))
	if rc != 0 {
		return helloif_getError(rc)
	}
	return nil
}
