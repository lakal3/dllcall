# Why DLLCall tool

Go already support embedding C/C++ code using CGO. So why is that not always suitable for C/C++ embendding?

## CGO challenges

CGO has some challenges especially in Windows

1\. CGO assumes GNU compatible compiler with archive file support
 
You cannot switch to a compiler that don't generate and consume unix style archive files. LIB files are standard on Windows.
Even newer CLANG compilers on Windows will consume and produce LIB files.

2\. CGO slows down go compilation

3\. Difficult binary distribution of C/C++ part
 
 Typically C/C++ interfaced module is existing project that is updated very seldom.
 With shared libraries you can precompile C/C++ module and distribute it in binary format and precompiled modules can be used even without installing C/C++ toolchain on every machine.   
 
 Windows don't have any C/C++ compiler installed by default and if it has, 
 it is most like Microsoft C/C++ compiler that currently can't be used for CGO programs.
 
 4\. Debugging tools
 
 Best Windows debuggers use PDB debug format. CGO toolchain cannot product PDB files for debugging in Windows.
 
 5\. Platform specific extensions
 
 For example MSVC supports some non standard extensions like importing typelibraries to interface with COM+ DLLs. 
 These can't be used from GNU compiler. 
 
 ## Raw dll calls
 
 In go, windows system calls are actually implemented using a mechanism that can load and
 call windows shared libraries (DLL). All Windows user level systems apis are implemented using a DLL.
 
 Raw calls to shared libraries may however be quite a challenging to set up properly if interface requires large structures. 
 Standard syscalls on windows have also some limitations like returning a floating point value from DLL call (used for example in sqlite3).
 
 Go syscall have some overhead for a good reason. To mitigate this overhead we can usually combine several raw calls into a single interface call. 
 
 ## Linux support
 
 Internally DLLCall uses go standard syscall mechanism on Windows to load shared libraries and invoke method from them. 
 DLLCall just simplifies interface generation and maintenance. 
 
 Unfortunately, go syscalls in Linux are implemented differently and 
 syscall mechanism cannot call Linux shared libraries (.so)
 
 It is possible to import dlopen/dlsym using CGO but this kind of ruins the original
 idea to support shared libraries without CGO. There is no pure go dlopen/dlsym implementation
 that I am aware of. See issue golang/go#18296 for more discussion about this.
 You can see hellolinux example on how to use interface with CGO this.
 
 Secondly, if we invoke loaded APIs through CGO, it will detect that we
 have pointers to go heap and will panic. Generally it is advisable not to send 
 pointer to go heap into C/C++ libraries as they are not aware of go's garbage collector.
 I this case interface code is aware of calling environment and passing go slices and strings should be ok. 
 
 If CGO could be used to load libraries we could still keep some of benefits like:
 - Ability to choose compiler (Clang) for shared library compilation
 - Faster compilation (we don't need to include actual project headers, only dllib)
 
 But currently this would require similar assembler code that is used to make Windows DLL calls 
   
   
 