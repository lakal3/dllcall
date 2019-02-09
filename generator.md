# DLLCall interface generation

## Scanning comments

DLLCall will use comments to attach metadata to types. This metadata will guide interface generation.
 
DLLCall tool will first scan source file for any comments attached to type definitions.

All comment lines starting with #c are reserved for generator. If you want to add comments about actual interface put those above any #c definitions.
Everything below first #c line will be copied to C++ interface. For example:


Currently we only support `#ctype` and `#cmethod`. If we find any other comment lines starting with #c generator will raise an error and abort generation.

### \#cmethod method_name
\#cmethod will define method name for interface. Only struct types may define methods.

Interface generator will generate named method for each \#cmethod definition. 
You can call there methods like any normal go methods. 

Each method will return an error and take no arguments. 
Pointer to structure is passed to DLL directly and all parameters and return values must be members of structure. 
Interface will have no separation for input or output values. We just pass pointer to a structure to DLL.

For example:

### \#ctype alias_definition

\#ctype will define type alias that we use in C++ for given go type. Typically we define some uintptr and 
and user real pointer type as alias in C++. For example

```go
/*
#ctype Stmt *
*/
type stmtHandle uintptr
```
Any field of type stmtHandle in go structure will be mapped to Stmt *.
**You must ensure that aliased type has same size in C++ and in go**

You can also define enums with type alias. For example:
```go
/*
Enum type for operand
#ctype operKind
enum operKind: int32_t {
  Get = 0,
  Put = 1,
  Delete = 2
};
 */
type OperKind int32
const (
	Get = OperKind(0)
	Put = OperKind(1)
	Delete = OperKind(2)
)
```

Anything after #c definitions will be copied to C++ interface.

## Running generator

After comment scanning generator will create a new go file that includes copy of interface
and additional code to generate actual interfaces. 

Generator will compile and run this code. It will use reflection to parse type information and 
use that information to generate actual C++ interface structures.

Generator temporary go file will be deleted after generation process. You can use -keep flag to preserve generator temporary go file. 
 
Currently generator supports following type that have well defined counterpart on C++ side.
- Any sized int and uint types (uint8, int8, uint16, int16, ...). 
They will be mapped equivalent C++ sized types (uint8_t, int8_t, ...)
- float32 and float64
- string. String is mapped to GoString and must be treated **readonly**. 
Writing to a string will most like cause program to crash.
- a pointer to supported types
- a slice of supported types
- an array of supported types
- Structures containing supported types

Due to copying of {interface}.go file, all types defined in this module that are used in interface types
must be defined in interface file. Otherwise generator will not find type definitions.
 
Not supported types includes for example maps and channels.

 
## Generated C++ code

For example, if we have definition:
```go
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

Generated C++ file will have a structure dbIf looking like:
```cpp
typedef   struct {
    Stmt * handle;
    GoString dbName;
    GoError *Open();
    GoError *Close();
} dbIf ;
```

For actual implementation we have to create method bodies for Open and Close methods.

 
