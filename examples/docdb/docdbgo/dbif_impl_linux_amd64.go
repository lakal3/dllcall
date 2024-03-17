package main

// Generated file. Do not edit

import "github.com/lakal3/dllcall/linux/syscall"
import "unsafe"
import "errors"
import "fmt"

var (
	_dbif_gate__getError  uintptr
	_dbif_gate_dbIf_Open  uintptr
	_dbif_gate_dbIf_Close uintptr
	_dbif_gate_dbBatch_Do uintptr
)

func load_dbif(dllPath string) (err error) {
	if _dbif_gate__getError != 0 {
		return nil
	}
	err = syscall.DisableCgocheck()
	if err != nil {
		return err
	}
	dll, err := syscall.LoadLibrary(dllPath)
	if err != nil {
		return err
	}
	_dbif_gate_dbIf_Open, err = syscall.GetProcAddress(dll, "dbIf_Open")
	if err != nil {
		return fmt.Errorf("%s: %v", "dbIf_Open", err)
	}
	_dbif_gate_dbIf_Close, err = syscall.GetProcAddress(dll, "dbIf_Close")
	if err != nil {
		return fmt.Errorf("%s: %v", "dbIf_Close", err)
	}
	_dbif_gate_dbBatch_Do, err = syscall.GetProcAddress(dll, "dbBatch_Do")
	if err != nil {
		return fmt.Errorf("%s: %v", "dbBatch_Do", err)
	}

	getcrc, err := syscall.GetProcAddress(dll, "GetCRC")
	if err != nil {
		return fmt.Errorf("GetCRC: %v", err)
	}
	var crc uint64
	syscall.SyscallN(getcrc, uintptr(unsafe.Pointer(&crc)))
	if crc != 0x98b5330a8380a2f0 {
		return fmt.Errorf("CRC mismatch %s != %x. DLL is not from same build than go code.", "0x98b5330a8380a2f0", crc)
	}
	_dbif_gate__getError, err = syscall.GetProcAddress(dll, "GetError")
	if err != nil {
		return err
	}
	return nil
}

func dbif_getError(rc uintptr) error {
	errText := make([]byte, 0, 512)
	syscall.SyscallN(_dbif_gate__getError, rc, uintptr(unsafe.Pointer(&errText)))
	return errors.New(string(errText))
}

func (r *dbIf) Open() (err error) {
	rc, _, _ := syscall.SyscallN(_dbif_gate_dbIf_Open, uintptr(unsafe.Pointer(r)), uintptr(24))
	if rc != 0 {
		return dbif_getError(rc)
	}
	return nil
}

func (r *dbIf) Close() (err error) {
	rc, _, _ := syscall.SyscallN(_dbif_gate_dbIf_Close, uintptr(unsafe.Pointer(r)), uintptr(24))
	if rc != 0 {
		return dbif_getError(rc)
	}
	return nil
}

func (r *dbBatch) Do() (err error) {
	rc, _, _ := syscall.SyscallN(_dbif_gate_dbBatch_Do, uintptr(unsafe.Pointer(r)), uintptr(24))
	if rc != 0 {
		return dbif_getError(rc)
	}
	return nil
}
