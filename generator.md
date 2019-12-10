# DLLCall interface generation

## Scanning comments

DLLCall will use comments to attach metadata to types. This metadata will guide interface generation.
 
DLLCall tool will first scan source file for any comments attached to type definitions.

All comment lines starting with #c are reserved for generator. If you want to add comments about the actual interface, put those above any #c definitions.
Everything below the first #c line will be copied into the C++ interface (.h-file).


Currently only `#ctype`, `#cmethod` and `#csafe_method` are supported. If any other comment is found, lines starting with #c generator will raise an error and abort generation.

### \#cmethod method_name
\#cmethod will define method name for the interface. Only struct types may define methods.

Interface generator will generate named method for each \#cmethod definition. 
Their methods can be called like any standard Go methods.

Each method will return an error and take no arguments. 
Pointer to structure is passed directly to DLL so all parameters and return values must be members of structure. 
Interface will have no separation for input or output values. We just pass a pointer to a structure to a DLL.


### \#ctype alias_definition

\#ctype will define the type alias that we use in C++ for the given Go type. Typically we define some uintptr 
and use a real pointer type as alias in C++. For example:

```go
/*
#ctype Stmt *
*/
type stmtHandle uintptr
```
Any field of type stmtHandle in a Go structure will be mapped to Stmt *.
**You must ensure that aliased type has exactly the same size in C++ and in go**

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

Anything after #c definitions will be copied to the C++ interface.

### \#csafe_method method_name (Experimental)

\#csafe_method works like method unless you also specify -fast option when running dllcall.
With fast options, dllcall will generate Go assembly file that will invoke method direcly, bypassing 
normal Go Syscall / CGO overhead. 

There are several limitations in fastcall implementation and it should be only used where perfomance is 
critical (for example calls to set OpenGl state). You can use new example fibon to experiment with actual call overhead.
In most cases overhead is not relevant.

To mark method as a safe_method it:
- MUST NOT use stack more that 63k. Exceeding this limit will crash your program!
- SHOULD NOT use any IO
- SHOULD NOT last over 1ms. Use normal call if in any doubt 
- CAN NOT return result(s) in structure variables. You may still use pointer to return values
as long as you are sure that these points to heap and not stack!

*In order to support Go assembler in fastcall interfaces, Linux DLL loader was moved to separate package linux/syscall. 
Go´s inbuilt assembler wont work if package has any CGO code. See issue [19948](https://github.com/golang/go/issues/19448) for details.*
 
This syscall library has similar methods to load and invoke functions from shared libraries that
already exists on Go´s standard syscall library on Windows.  

## Running generator

After comment scanning, the generator will create a new Go file that includes a copy of the interface
and additional code to generate actual interfaces. 

Generator will compile and run this code. It will use reflection to parse type information and 
use that information to generate actual C++ interface file.

Generator temporary Go file will be deleted after generation process. You can use -keep flag to preserve generator temporary go file. 
 
Currently generator supports following types that have well defined counterparts on C++ side.
- Any size of int and uint types (uint8, int8, uint16, int16, ...). 
They will be mapped to C++ equivalents (uint8_t, int8_t, ...)
- float32 and float64
- String. String is mapped to GoString (see Go types) 
- a pointer to supported types
- a slice of supported types (see Go types)
- an array of supported types
- Structures containing supported types

Due to the copying of {interface}.go file, all types defined in this module that are used in interface types
must be defined in interface file. Otherwise generator will not find type definitions.
 
Unsupported types include for example maps and channels.

 
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

The generated C++ file will have a structure dbIf:
```cpp
typedef   struct {
    Stmt * handle;
    GoString dbName;
    GoError *Open();
    GoError *Close();
} dbIf ;
```

For actual implementation we have to create method bodies for Open- and Close-methods.

 
## Go types

Generated C++ header file will include default implementations of certain Go types that
matches Go memory layout:
- GoString - Matches Go's string definition. **GoString are always utf-8, readonly and they don't have terminating 0 character.**
Use append method to copy Go's string to std::string before you manipulate it.

- GoSlice<T> - Go slice of type. You can safely change length and content within capacity boudaries.

- GoError - Support to return error values back to Go code. Each generated function
   will return a GoError struct on nullptr to indicate that there was no error. Go interface will call
   GoError::GetError to retrieve actual error message into byte slice.

You can override default implementation by defining your own type and declaring DLLCALL_CUSTOM_GO_XXX. 
See generated interface file for more details.
 