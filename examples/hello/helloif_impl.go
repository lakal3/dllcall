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
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_helloif_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
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
	crc, _, _ := syscall.Syscall(getcrc, 0, 0, 0, 0)
	if uint64(crc) != 0x00715e677960c0ae {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x00715e677960c0ae", crc)
	}
	return nil
}

func getError_helloif(rc uintptr) error {
	errText := make([]byte, 0, 512)
	syscall.Syscall(_helloif_gate__getError, 2, rc, uintptr(unsafe.Pointer(&errText)), 0)
	return errors.New(string(errText))
}

func (r *greeting) Greet() (err error) {
	rc, _, _ := syscall.Syscall(_helloif_gate_greeting_Greet, 2, uintptr(unsafe.Pointer(r)),
		uintptr(16), 0)
	if rc != 0 {
		return getError_helloif(rc)
	}
	return nil
}
