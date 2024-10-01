package main

// Generated file. Do not edit

import "github.com/lakal3/dllcall/linux/syscall"
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
	syscall.GetCRC(getcrc, &crc)
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
	errText := make([]byte, 0, 4096)
	syscall.GetError(_helloif_gate__getError, rc, &errText)
	return errors.New(string(errText))
}

func (r *greeting) Greet() (err error) {
	rc := syscall.SyscallT(_helloif_gate_greeting_Greet, r)
	if rc != 0 {
		return helloif_getError(rc)
	}
	return nil
}
