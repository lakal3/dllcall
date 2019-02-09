package main

/*
#cgo LDFLAGS: -ldl
#include <dlfcn.h>

void *invoke(void *addr, void *arg, long long size) {
	void *(*ptr)(void *arg, long long size) = addr;
	return (*ptr)(arg, size);
}
*/
import "C"
import (
	"log"
	"unsafe"
)

func main() {
	handle := C.dlopen(C.CString("./libgreeting.so"), C.RTLD_NOW)
	if uintptr(handle) == 0 {
		log.Fatal("Load libgreeting.so failed")
	}
	addr := C.dlsym(handle, C.CString("greeting_Greet"))
	if uintptr(addr) == 0 {
		log.Fatal("Load libgreeting.so failed")
	}
	// FAILS: gr := &greeting{ greeting: fmt.Sprintf("Hello DLLCall!. Time is %v", time.Now())}
	// panic: runtime error: cgo argument has Go pointer to Go pointer
	gr := &greeting{greeting: "Hello DLLCall!"}
	C.invoke(addr, unsafe.Pointer(gr), C.longlong(unsafe.Sizeof(greeting{})))

}

type greeting struct {
	greeting string
}
