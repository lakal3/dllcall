# DLLCall tool

DLLCall is an iterface generator. 

DLLCall will use single go file as an interface desctiption and
from that file it will generate wrappers to call DLL implementing same interface.

See [why DLLCall](why.md) for comparision between CGO interface and DLLCall interfaces.

## Using DLLCall

Install DLLCall tool using `go install github.com/lakal2/dllcall`

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
When we invoke dllcall tool like `dllcall dbif.go {otherdir}/dbif.h`, 
dllcall will generate dbif_impl.go file into current directory and dbif.h to specifier location.

dbif_impl.go file name is derived from original go file name. C++ file can be freely named.

You can use go generate directive to and go generate tool to easily upgrade interface whenever it changes.

## Interface definition

Interface definition must be a single go file containing only type definitions. 
You can't add functions or methods to interface file.

## Interface usage
After you have generated interface, you can call it like any other go 

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
Generator will create {interface}_impl.go file that contains add methods 
defined with #cmethod. 

It will also contain load\_{interface} method that takes a single file path to shared library (dll).
How libraries are located depend on your operating system. Your must call load\_{interface} before you try to invoke any other methods from interface.

## Interface implementation

Generator will create .h file that you must include into shared library project. 

You must also implement generated method stubs in shared library project. 
In this example methods to implement are:
```cpp
GoError *dbIf::Open();
GoError *dbIf::Close();
```

For more detailed description of comment annotations, generation process and supported types see [generator](generator.md)
or examples.

# Status
This project is still in early stage but has been successfully user to embed some
notable C/C++ libraries including sqlite3, SDL2 and several Window only COM+ programs.

Currently only 64 bit (amd64) Windows in supported. 

Breaking changes are still possible but not very likely.

## TODO?
- [ ] Sqlite example project using DLLCall
- [ ] Linux support - See [why DLLCall](why.md) for actual challenges with this
- [ ] Win32 support




  


