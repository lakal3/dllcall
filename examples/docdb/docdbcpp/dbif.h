
// Generated file. Do not edit
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
    static void GetError(GoError *err, GoSlice<char> &errBuf) {
        size_t len = err->error.size();
        if (len >= errBuf.cap) { len = errBuf.cap - 1; }
        strncpy(errBuf.data, &(err->error.at(0)), len);
        errBuf.len = len;
        delete err;
    }

};
#endif





enum operKind: int32_t {
  Get = 0,
  Put = 1,
  Delete = 2
};



typedef   struct dbIf {
    Stmt * handle;
    GoString dbName;
    GoError *Open();
    GoError *Close();
} dbIf ;
typedef   struct dbOper {
    operKind kind;
    GoString key;
    GoSlice<uint8_t > value;
} dbOper ;
typedef   struct dbBatch {
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
DLL_EXPORT void DLLCALL_SYSCALL GetError(GoError *err, GoSlice<char> *errBuf);
DLL_EXPORT void DLLCALL_SYSCALL GetCRC(uint64_t *crc);
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

void GetError(GoError *err, GoSlice<char> *errBuf) {
	return GoError::GetError(err, *errBuf);
}

void GetCRC(uint64_t *crc) {
    *crc = 0x98b5330a8380a2f0ull;
}
#endif
