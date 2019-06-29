
// Generated file. Not not edit
#define _CRT_SECURE_NO_WARNINGS

#include <cstdint>
#include <string>
#include <cstring>

#ifndef DLLCALL_CUSTOM_GO_SLICE
#define DLLCALL_CUSTOM_GO_SLICE
template<class T> 
struct GoSlice
{
    T *data;
    uint64_t len;
    uint64_t cap;
};
#endif

#ifndef DLLCALL_CUSTOM_GO_STRING
#define DLLCALL_CUSTOM_GO_STRING
struct GoString
{
    const char *data;
    uint64_t len;
    void append(std::string &str) { str.append(data, len); }
};
#endif

#ifndef DLLCALL_CUSTOM_GO_ERROR
#define DLLCALL_CUSTOM_GO_ERROR
struct GoError {
    std::string error;
    GoError(const char *err): error(err) {
    }
    const char *GetError() { return error.c_str(); }
    static void GetError(GoError *err, GoSlice<char> errBuf) {
        size_t len = strlen(err->GetError());
        if (len >= errBuf.cap) { len = errBuf.cap - 1; }
        strncpy(errBuf.data, err->GetError(), len);
        errBuf.len = len;
        delete err;
    }

};
#endif


typedef   struct greeting {
    GoString greeting;
    GoError *Greet();
} greeting ;
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
DLL_EXPORT void DLLCALL_SYSCALL GetCRC(uint64_t *crc);
DLL_EXPORT GoError * DLLCALL_SYSCALL greeting_Greet(greeting *arg, int64_t argLen );
}
#ifndef DLLCALL_NO_IMPL
const char *_callError = "Argument length check failed. Recompile interface and check compiler alignments";
GoError *greeting_Greet(greeting *arg, int64_t argLen ) {
    if (sizeof(greeting) != argLen) { return new GoError(_callError); }
    GoError *err;
    err = arg->Greet();
    return err;
}

void GetError(GoError *err, GoSlice<char> errBuf) {
	return GoError::GetError(err, errBuf);
}

void GetCRC(uint64_t *crc) {
    *crc = 0x21e8f5462f9aeab2ull;
}
#endif
