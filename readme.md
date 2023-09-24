# DLLCall tool

DLLCall is an interface generator. 

DLLCall will use a single go file as an interface description and
from that file it will generate wrappers to call DLL implementing the same interface.

See [why DLLCall](why.md) for comparision between CGO and DLLCall interfaces.

## Using DLLCall

Install DLLCall tool using `go install github.com/lakal3/dllcall`

Invoke:  `dllcall {interface}.go {c++ interface}.h`

# Example interface dbif.go

```go
package {package name}

/*
#ctype Stmt *
*/
type stmtHandle uintptr

/*
#cmethod Open
#cmethod Close
*/
type dbIf struct{
	handle stmtHandle
	dbName string	
}
```

## Generate interface
When we invoke dllcall tool such as `dllcall dbif.go {otherdir}/dbif.h`, 
dllcall will generate a dbif_impl_windows_amd64.go and dbif_impl_linux_amd64.go file into the current directory and dbif.h to a specified location.

The file name of dbif_impl_{os}_amd64.go is derived from the original go file name. C++ -file can be freely named.

You can use `//go:generate` directive and `go generate` command to easily upgrade the interface whenever it changes.

## Interface definition

Interface definition must be a single go file containing only type definitions. 
You can't add functions or methods to the interface file.

## Interface usage
After you have generated the interface, you can call it like any other go method 

```go
func main()  {
	err := load_dbif("docdbdll.dll") 
	// on Linux:  err := load_dbif(".\libdocdbdll.so") 
	if err != nil { log.Fatal( err )}
	d := &dbIf{ dbName: "test.db"}
	err = d.Open()
	if err != nil { log.Fatal( err )}
	/// ... Do something with db
	err = d.Close()
	if err != nil { log.Fatal( err )}
}
```

Generator will create an {interface}\_impl\_{os}_amd64.go file that contains all methods 
defined with #cmethod. 

It will also contain load\_{interface} method that takes a single file path to shared library (dll).
How libraries are located depends on your operating system.
 
You must call load\_{interface} before you try to invoke any other methods from interface.

## Interface implementation

Generator will create a .h-file that you must include into the shared library project. 

You must also implement the generated method stubs in a shared library project. 
In this example the methods to implement are:
```cpp
GoError *dbIf::Open();
GoError *dbIf::Close();
```

For more detailed description of comment annotations, generation process and supported types, see [generator](generator.md)
or examples.

## Cgocheck

Go 1.16 introduced new check that prevents pointers structures that contains pointer to Go memory.
But in dllcall library all code must be aware of Go garbage collector and do not retain any references
to structure member after call has been completed, so this check is unnecessary.

There are three ways to suppress this check

### Environment variable

Add environment variable GODEBUG=cgocheck=0. This is awkward as we must set flag before running any executable using dllcalls.

### Pinned memory
In version 1.21 Go introduced new Pin mechanism that allows marking all pointer that we want to use in calls to C++ program so
that Gos garbage collector is aware of them. 

You can enable generating pinned memory pointer with -pin flag. 

Pinning is a standard way to handle pointers to Go memory, but unfortunately it add some overhead to calls. Hello Example uses this method.

### Disable cgocheck with code

Dllcall uses hack to turn off cgocheck. This hack is in function DisableCgocheck. 
Generated code will call this function when loading .dll/.so file.

If you select -pin option, cgocheck is not disabled because it will allow pinned memory pointers without panic.

It is possible this hack will not work in some future version of Go, so you can change this function
to disable this hack. 

### Windows

Actually Windows syscalls don't apply any checks for go pointers, so dllcall generator will
not emit any code to disable checks of to pin memory even if -pin flag is given


## Safe method (experimental)

*In go 1.21.1 safemethod calls will fail because they set new value to SP (Stack pointer). 
This is due to change in go1.21 witch will hopefully be fixed in go 1.21.2*

Added new experimental \#csafe_method that mostly allows bypassing Go call overhead to native libraries.
This method has several limitations and is experimental. See [Generator](generator.md) and new sample fibon for more details.


# Status
This project is still in early stages, but has been successfully used to embed some
notable C/C++ libraries including sqlite3, SDL2 and several Windows only COM+ programs.

Currently only 64 bit (amd64) Windows and 64 bit Linux (amd64) are supported. 

Breaking changes are still possible but not very likely.

**Version 0.8.2 breaking changes**
- Renamed generated Go interface. Generated files have now _amd64 extension to prevent build on other architectures. 
To fix this, remove old generated files from project.
- Linux CGO wrapper has been moved to a new package linux/syscall to support fastcall. Wrapper packages must be compiled without CGO.

   

## TODO?
- [x] Pinned memory available with Go 1.21
- [x] Fastcall (Experimental)
- [ ] Better support for types imported from other modules
- [x] Linux support - Implemented but have some issues (see [Why DLL call](why.md))
- [ ] Sqlite example project using DLLCall
- [ ] 32 bit support 






  


