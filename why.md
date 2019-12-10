# Why DLLCall tool

Go already supports embedding C/C++ code using CGO. So why is CGO not always suitable for C/C++ embedding?

## CGO challenges

CGO has some challenges, especially in Windows

1\. CGO assumes GNU compatible compiler with archive file support
 
You cannot switch to a compiler that does not generate and consume unix style archive files. LIB files are standard on Windows.
Even newer CLANG compilers on Windows will consume and produce LIB files.

2\. CGO slows down Go compilation

3\. Difficult binary distribution of C/C++ part
 
 Typically C/C++ interfaced module is an existing project that is updated very seldom.
 With shared libraries you can precompile C/C++ module and distribute it in binary format. 
 Precompiled modules can be used even without installing C/C++ toolchain on every machine.   
 
 By default, Windows does not have any C/C++ compilers installed and if it has, 
 it is most likely a Microsoft C/C++ compiler which currently can't be used for CGO programs.
 
 4\. Debugging tools
 
 Best Windows debuggers use PDB debug format. CGO toolchain cannot produce PDB files for debugging in Windows.
 
 5\. Platform specific extensions
 
 For example MSVC supports some non standard extensions like importing typelibraries to use with COM+ DLLs. 
 This can't be done with GNU compiler. 
 
 ## Raw dll calls
 
 In go, the system calls in Windows are actually implemented using a mechanism that can load and
 call shared libraries (DLL). All user level apis are implemented using DLLs.
 
 However, raw calls to shared libraries may be quite challenging to set up properly if an API call requires large structures. 
 Standard syscalls on Windows have also some limitations like returning a floating point value from a DLL call (used for example in sqlite3).
 
 Go syscalls have some overhead for a good reason. To mitigate this overhead we can usually combine several raw calls into a single interface call. 
 
 ## Linux support
 
 Internally DLLCall uses Go standard syscall mechanism on Windows to load shared libraries and invoke methods from them. 
 DLLCall just simplifies interface generation and maintenance. 
 
 Unfortunately, syscalls in Linux are implemented differently and 
 syscall mechanism cannot call Linux shared libraries (.so). Dllcall has package linux/syscall that simulates
 Windows syscall functionalities. There methods are implemented with CGO and dlopen/dlsym methods.
 
 *There is no pure Go dlopen/dlsym implementation
 that I am aware of. See issue golang/go#18296 for more discussion about this.*
 
 In addition, if we invoke loaded APIs through CGO, it will detect that we
 have pointers to Go heap and will panic. **You can and must disable this check by setting environment variable GODEBUG=cgocheck=0** 
 
 Generally it is advisable not to send a pointer to Go-heap into C/C++ libraries as they are not aware of Go's garbage collector.
 In this case the interface code is aware of calling the environment so passing Go slices and strings should be ok.
 
 Even when CGO is used to load libraries we can keep certain benefits such as:
 - Ability to choose compiler (Clang) for shared library compilation
 - Faster compilation (we don't need to include actual project headers, only dllib)
  
 
   
   
 