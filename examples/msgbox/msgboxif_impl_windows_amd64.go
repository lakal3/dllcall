package main

// Generated file. Do not edit

import "syscall"
import "unsafe"
import "errors"
import "fmt"

var (
	_msgboxif_gate__getError   uintptr
	_msgboxif_gate_msgBox_show uintptr
)

func load_msgboxif(dllPath string) (err error) {
	if _msgboxif_gate__getError != 0 {
		return nil
	}
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_msgboxif_gate_msgBox_show, err = syscall.GetProcAddress(dll, "msgBox_show")
	if err != nil {
		return fmt.Errorf("%s: %v", "msgBox_show", err)
	}

	getcrc, err := syscall.GetProcAddress(dll, "GetCRC")
	if err != nil {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	_, _, _ = syscall.SyscallN(getcrc, uintptr(unsafe.Pointer(&crc)))
	if crc != 0x63a150b21b59a5fc {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x63a150b21b59a5fc", crc)
	}
	_msgboxif_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
	if err != nil {
		return err
	}
	return nil
}

func msgboxif_getError(rc uintptr) (err error) {
	errText := make([]byte, 0, 512)
	_, _, _ = syscall.SyscallN(_msgboxif_gate__getError, rc, uintptr(unsafe.Pointer(&errText)))
	return errors.New(string(errText))
}

func (r *msgBox) show() (err error) {
	rc, _, _ := syscall.SyscallN(_msgboxif_gate_msgBox_show, uintptr(unsafe.Pointer(r)), uintptr(32))
	if rc != 0 {
		return msgboxif_getError(rc)
	}
	return nil
}
