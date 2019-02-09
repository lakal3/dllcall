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
dllcall will generate a dbif_impl.go file into the current directory and dbif.h to a specified location.

The file name of dbif_impl.go is derived from the original go file name. C++ -file can be freely named.

You can use `//go:generate` directive and `go generate` command to easily upgrade the interface whenever it changes.

## Interface definition

Interface definition must be a single go file containing only type definitions. 
You can't add functions or methods to the interface file.

## Interface usage
After you have generated the interface, you can call it like any other go method 

```go
func main()  {
	err := load_dbif("docdbdll.dll")
	if err != nil { log.Fatal( err )}
	d := &dbIf{ dbName: "test.db"}
	err = d.Open()
	if err != nil { log.Fatal( err )}
	/// ... Do something with db
	err = d.Close()
	if err != nil { log.Fatal( err )}
}
```
Generator will create an {interface}_impl.go file that contains all methods 
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

# Status
This project is still in early stages, but has been successfully used to embed some
notable C/C++ libraries including sqlite3, SDL2 and several Windows only COM+ programs.

Currently only 64 bit (amd64) Windows is supported. 

Breaking changes are still possible but not very likely.

## TODO?
- [ ] Sqlite example project using DLLCall
- [ ] Linux support - See [why DLLCall](why.md) for actual challenges within this
- [ ] Win32 support




  


