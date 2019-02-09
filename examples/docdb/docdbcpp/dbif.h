
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





enum operKind: int32_t {
  Get = 0,
  Put = 1,
  Delete = 2
};



typedef   struct {
    Stmt * handle;
    GoString dbName;
    GoError *Open();
    GoError *Close();
} dbIf ;
typedef   struct {
    operKind kind;
    GoString key;
    GoSlice<uint8_t > value;
} dbOper ;
typedef   struct {
    GoSlice<dbOper > operations;
    GoError *Do();
} dbBatch ;
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
DLL_EXPORT GoError * DLLCALL_SYSCALL dbIf_Open(dbIf *arg, int64_t argLen );
DLL_EXPORT GoError * DLLCALL_SYSCALL dbIf_Close(dbIf *arg, int64_t argLen );
DLL_EXPORT GoError * DLLCALL_SYSCALL dbBatch_Do(dbBatch *arg, int64_t argLen );
}
#ifndef DLLCALL_NO_IMPL
const char *_callError = "Argument length check failed. Recompile interface and check compiler alignments";
GoError *dbIf_Open(dbIf *arg, int64_t argLen ) {
    if (sizeof(dbIf) != argLen) { return new GoError(_callError); }
    GoError *err;
    err = arg->Open();
    return err;
}
GoError *dbIf_Close(dbIf *arg, int64_t argLen ) {
    if (sizeof(dbIf) != argLen) { return new GoError(_callError); }
    GoError *err;
    err = arg->Close();
    return err;
}
GoError *dbBatch_Do(dbBatch *arg, int64_t argLen ) {
    if (sizeof(dbBatch) != argLen) { return new GoError(_callError); }
    GoError *err;
    err = arg->Do();
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
    return 0x98b5330a8380a2f0;
}
#endif
