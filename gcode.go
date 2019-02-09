package main

const generatorCode = `
func main() {
	f, err := os.Create("{{ .TargetFile }}")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	fgo, err := os.Create("{{ .GoTargetFile }}")
	if err != nil {
		log.Fatal(err)
	}
	defer fgo.Close()

	fmt.Fprintln(f, "{{ Quote .Header }}")
	{{ range .CDecl }}
	fmt.Fprintln(f, "{{ Quote .}}") 
{{ end }}
	var s string
	{{ range .GoTypes }}
	s = _gen_CType("{{ . }}", reflect.TypeOf(new({{.}})).Elem())
	fmt.Fprintln(f, "typedef ", s, ";");
{{ end }}
	fmt.Fprintln(f, "#ifndef DLL_EXPORT")
    fmt.Fprintln(f, "#ifdef _WIN32")
	fmt.Fprintln(f, "#define DLL_EXPORT  __declspec(dllexport) ")
    fmt.Fprintln(f, "#define DLLCALL_SYSCALL  __stdcall")
    fmt.Fprintln(f, "#else")
	fmt.Fprintln(f, "#define DLL_EXPORT")
    fmt.Fprintln(f, "#define DLLCALL_SYSCALL")
	fmt.Fprintln(f, "#endif")
	fmt.Fprintln(f, "#endif")
    fmt.Fprintln(f, "extern \"C\" {")
	fmt.Fprintln(f, "DLL_EXPORT void DLLCALL_SYSCALL GetError(GoError *err, GoSlice<char> errBuf);")
    fmt.Fprintln(f, "DLL_EXPORT uint64_t DLLCALL_SYSCALL GetCRC();")
	{{ range .Methods }}
	fmt.Fprintln(f, "DLL_EXPORT GoError * DLLCALL_SYSCALL {{ .GoType}}_{{ .MethodName }}({{ .GoType}} *arg, int64_t argLen );")
{{ end }}
	fmt.Fprintln(f, "}")	
	fmt.Fprintln(f, "#ifndef DLLCALL_NO_IMPL")
	fmt.Fprintln(f, "const char *_callError = \"Argument length check failed. Recompile interface and check compiler alignments\";")
{{ range .Methods }}
	fmt.Fprintln(f, "GoError *{{ .GoType}}_{{ .MethodName }}({{ .GoType}} *arg, int64_t argLen ) {")
	fmt.Fprintln(f, "    if (sizeof({{ .GoType}}) != argLen) { return new GoError(_callError); }")
	fmt.Fprintln(f, "    GoError *err;")
	fmt.Fprintln(f, "    err = arg->{{ .MethodName }}();")
	fmt.Fprintln(f, "    return err;")
	fmt.Fprintln(f, "}")
{{ end }}
    fmt.Fprintln(f, "#ifndef DLLCALL_CUSTOM_GO_ERROR")
	fmt.Fprintln(f, "void GetError(GoError *err, GoSlice<char> errBuf) {")
	fmt.Fprintln(f, "    size_t len = strlen(err->GetError());");
	fmt.Fprintln(f, "    if (len >= errBuf.cap) { len = errBuf.cap - 1; }")
	fmt.Fprintln(f, "	 strncpy(errBuf.data, err->GetError(), len);")
	fmt.Fprintln(f, "    errBuf.len = len;")
	fmt.Fprintln(f, "    delete err;")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "#endif")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "uint64_t GetCRC() {")
	fmt.Fprintf(f, "    return {{ .CRC }};\n");
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "#endif")

	// Go part
	fmt.Fprintln(fgo, "package {{ .PackageName }}")
	fmt.Fprintln(fgo, "// Generated file. Not not edit")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "import \"syscall\"")
	fmt.Fprintln(fgo, "import \"unsafe\"")
	fmt.Fprintln(fgo, "import \"errors\"")
	fmt.Fprintln(fgo, "import \"fmt\"")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "var (")
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate__getError uintptr")
{{ range .Methods }}
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }} uintptr")
{{ end }}
	fmt.Fprintln(fgo, ")")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func load_{{ .ModuleName }}(dllPath string)(err error) {")
	fmt.Fprintln(fgo, "	   dll, err := syscall.LoadLibrary(dllPath)")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate__getError, err = syscall.GetProcAddress(dll, \"GetError\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
{{ range .Methods }}
	fmt.Fprintln(fgo, "     _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, err = syscall.GetProcAddress(dll, \"{{ .GoType}}_{{ .MethodName }}\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return fmt.Errorf(\"%s: %v\", \"{{ .GoType}}_{{ .MethodName }}\",err)")
	fmt.Fprintln(fgo, "    }")
{{ end }}
    fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "    getcrc, err := syscall.GetProcAddress(dll,\"GetCRC\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return fmt.Errorf(\"GetCRC: %v\", err)")
	fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "    crc, _, _ := syscall.Syscall(getcrc,0,0,0,0)")
    fmt.Fprintln(fgo, "    if uint64(crc) != {{ .CRC }} {")
    fmt.Fprintln(fgo, "        return fmt.Errorf(\"CRC mismatch %s != %x. DLL is not from same build than go code.\",\"{{ .CRC }}\", crc)")
    fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "	   return nil")
	fmt.Fprintln(fgo, "}")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func getError_{{ $.ModuleName }}(rc uintptr)(error) {")
	fmt.Fprintln(fgo, "    errText := make([]byte, 0, 512)")
	fmt.Fprintln(fgo, "    syscall.Syscall(_{{ $.ModuleName }}_gate__getError, 2, rc, uintptr(unsafe.Pointer(&errText)), 0)")
	fmt.Fprintln(fgo, "    return errors.New(string(errText))")
	fmt.Fprintln(fgo, "}")
	fmt.Fprintln(fgo, "")
{{ range .Methods }}
	fmt.Fprintln(fgo, "func (r *{{ .GoType}}) {{ .MethodName }}() (err error) {")
	fmt.Fprintln(fgo, "    rc, _, _ := syscall.Syscall(_{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, 2, uintptr(unsafe.Pointer(r)), ")
    fmt.Fprintln(fgo, "      uintptr(", reflect.TypeOf(new({{ .GoType}})).Elem().Size(), "), 0)")
    fmt.Fprintln(fgo, "    if rc != 0 {")
    fmt.Fprintln(fgo, "         return getError_{{ $.ModuleName }}(rc)")
    fmt.Fprintln(fgo, "    }")
    fmt.Fprintln(fgo, "    return nil")
	fmt.Fprintln(fgo, "}")
    fmt.Fprintln(fgo, "")
{{ end }}
}

func _gen_CType(typeName string, t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		tn := _gen_CType("", t.Elem())
		return "(" + tn + "*) " + typeName
	}
	if t.Kind() == reflect.Slice {
		tn := _gen_CType("", t.Elem())
		return "GoSlice<" + tn + "> " + typeName
	}
	alias, ok := _gen_alias[t.Name()]
	if ok {
		return alias + " " + typeName;
	}

	if t.Kind() == reflect.Array {
		tn := _gen_CType("",t.Elem())
		_gen_alias[t.Name()] = t.Name()
		return fmt.Sprintf("%s %s[%d]", tn, typeName, t.Len()) 
	}

	if t.Kind() == reflect.Struct {
		sb := &strings.Builder{}
		sb.WriteString(" struct {\n")
		for idx := 0; idx < t.NumField(); idx++ {
			f := t.Field(idx)
			tn := _gen_CType(f.Name, f.Type)
			sb.WriteString("    " + tn + ";\n")
		}

{{ range .Methods }}
		if typeName == "{{ .GoType }}" {
			sb.WriteString("    GoError *{{.MethodName }}();\n")
		}
{{ end }}
		sb.WriteString("}")
		sb.WriteString(" ")
		sb.WriteString(typeName)
		_gen_alias[t.Name()] = t.Name()
		return sb.String()
	}
	return t.Name()
}

var _gen_alias = map[string]string {
{{ range .Aliases}}    "{{ .GoType }}" : "{{ .CAlias }}",
{{ end }}
}
`

const header = `
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

`
