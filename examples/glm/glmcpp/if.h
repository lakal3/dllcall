
// Generated file. Not not edit
#define _CRT_SECURE_NO_WARNINGS

#include <cstdint>
#include <string>
#include <cstring>

template<class T> 
struct GoSlice
{
    T *data;
    uint64_t len;
    uint64_t cap;
};

struct GoString
{
    const char *data;
    uint64_t len;
    void append(std::string &str) { str.append(data, len); }
};

#ifndef DLLCALL_CUSTOM_GO_ERROR
struct GoError {
    std::string error;
    GoError(const char *err): error(err) {}
    const char *GetError() { return error.c_str(); }
};
#endif


typedef  double  Vector[3] ;
typedef  double  Matrix4x4[16] ;
typedef   struct {
    Matrix4x4 Mat;
    GoSlice<Vector > Vectors;
    GoError *Multiply();
} MultiplyVectors ;
#ifndef DLL_EXPORT
#ifdef _WIN32
#define DLL_EXPORT  __declspec(dllexport) 
#define DLLCALL_SYSCALL  __stdcall
#else
#define DLL_EXPORT
#define DLLCALL_SYSCALL
#endif
#endif
extern "C" {
DLL_EXPORT void DLLCALL_SYSCALL GetError(GoError *err, GoSlice<char> errBuf);
DLL_EXPORT uint64_t DLLCALL_SYSCALL GetCRC();
DLL_EXPORT GoError * DLLCALL_SYSCALL MultiplyVectors_Multiply(MultiplyVectors *arg, int64_t argLen );
}
#ifndef DLLCALL_NO_IMPL
const char *_callError = "Argument length check failed. Recompile interface and check compiler alignments";
GoError *MultiplyVectors_Multiply(MultiplyVectors *arg, int64_t argLen ) {
    if (sizeof(MultiplyVectors) != argLen) { return new GoError(_callError); }
    GoError *err;
    err = arg->Multiply();
    return err;
}
#ifndef DLLCALL_CUSTOM_GO_ERROR
void GetError(GoError *err, GoSlice<char> errBuf) {
    size_t len = strlen(err->GetError());
    if (len >= errBuf.cap) { len = errBuf.cap - 1; }
	 strncpy(errBuf.data, err->GetError(), len);
    errBuf.len = len;
    delete err;
}
#endif

uint64_t GetCRC() {
    return 0x00698cefaddc526c;
}
#endif
