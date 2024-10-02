# DLLCall tool

DLLCall is an interface generator. 

DLLCall will use a single go file as an interface description and
from that file it will generate wrappers to call DLL implementing the same interface.

See [why DLLCall](why.md) for comparison between CGO and DLLCall interfaces.

## Using DLLCall

Install DLLCall tool using `go install github.com/lakal3/dllcall`. 

You must have at least Go 1.18 to use DLLCall

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
Dllcall will bypass cgocheck, because generated interface code don't have full type information to support this.

When using dllcall all code must be aware of Go garbage collector and do not retain any references
to structure member in C++ side after call has been completed. If you want to keep go allocated heap elements in C++ code you can use new runtime.Pinner available from Go 1.21.



## Safe method (experimental)

Added new experimental #csafe_method that mostly allows bypassing Go call overhead to native libraries.
This method has several limitations and is experimental. See [Generator](generator.md) and new sample fibon for more details.

You should only use fast call when absolutely necessary like accessing high resolution timer in Windows. 


*In go 1.21.1 safemethod calls will fail because they set new value to SP (Stack pointer).
This is due to change in go1.21. It has been fixed in go1.22 for Windows but Linux still.
It seems to work in go1.23 for both operating systems*

# Status

Dllcall is very near release version 1.0. If 0.12.1 don't have any issues it will promoted to version 1.0.0.

Currently only 64 bit (amd64) Windows and 64 bit Linux (amd64) are supported. 
 

## TODO?
- [x] Fastcall (Experimental)
- [ ] Better support for types imported from other modules
- [x] Linux support 
- [ ] Pin example






  


