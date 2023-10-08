package main

const generatorCode1 = `
func main() {
	__buildHFile()
{{ if .BuildWindows }}
	__winLoader()
{{     if .SafeMethods }}
	__winFastcall()
{{     end }}
{{ end }}
{{ if .BuildLinux }}
	__linuxLoader()
{{     if .SafeMethods }}
    __linuxFastcall()
{{     end }}
{{ end }}
}

func __buildHFile() {
	f, err := os.Create("{{ .TargetFile }}")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	
	fmt.Fprintln(f, "{{ Quote .Header }}")
	{{ range .CDecl }}
	fmt.Fprintln(f, "{{ Quote .}}") 
{{ end }}
	var s string
	{{ range .GoTypes }}
	s = _gen_CType("{{ NoModule . }}", reflect.TypeOf(new({{.}})).Elem())
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
	fmt.Fprintln(f, "DLL_EXPORT void DLLCALL_SYSCALL GetError(GoError *err, GoSlice<char> *errBuf);")
    fmt.Fprintln(f, "DLL_EXPORT void DLLCALL_SYSCALL GetCRC(uint64_t *crc);")
{{ range .Methods }}
	fmt.Fprintln(f, "DLL_EXPORT GoError * DLLCALL_SYSCALL {{ .GoType}}_{{ .MethodName }}({{ .GoType}} *arg, int64_t argLen );")
{{ end }}
{{ range .SafeMethods }}
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
{{ range .SafeMethods }}
	fmt.Fprintln(f, "GoError *{{ .GoType}}_{{ .MethodName }}({{ .GoType}} *arg, int64_t argLen ) {")
	fmt.Fprintln(f, "    if (sizeof({{ .GoType}}) != argLen) { return new GoError(_callError); }")
	fmt.Fprintln(f, "    GoError *err;")
	fmt.Fprintln(f, "    err = arg->{{ .MethodName }}();")
	fmt.Fprintln(f, "    return err;")
	fmt.Fprintln(f, "}")
{{ end }}
    fmt.Fprintln(f, "")
	fmt.Fprintln(f, "void GetError(GoError *err, GoSlice<char> *errBuf) {")
	fmt.Fprintln(f, "	return GoError::GetError(err, *errBuf);")
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "")
	fmt.Fprintln(f, "void GetCRC(uint64_t *crc) {")
	fmt.Fprintln(f, "    *crc = {{ .CRC }}ull;");
	fmt.Fprintln(f, "}")
	fmt.Fprintln(f, "#endif")
}
`

const loadWin = `
func __winLoader() {
	fgo, err := os.Create("{{ .GoTargetFile }}_windows_amd64.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fgo.Close()
	
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "//")
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
{{ range .SafeMethods }}
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }} uintptr")
{{ end }}
	fmt.Fprintln(fgo, ")")
{{ if .SafeMethods }}
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func _{{ .ModuleName }}_fastcall(trap uintptr, ptr uintptr, size uintptr)(errPtr uintptr)")
{{ end }}
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func load_{{ .ModuleName }}(dllPath string)(err error) {")
	fmt.Fprintln(fgo, "    if _{{ $.ModuleName }}_gate__getError != 0 {")
	fmt.Fprintln(fgo, "        return nil")
	fmt.Fprintln(fgo, "    }")

	fmt.Fprintln(fgo, "    dll, err := syscall.LoadLibrary(dllPath)")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
{{ range .Methods }}
	fmt.Fprintln(fgo, "     _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, err = syscall.GetProcAddress(dll, \"{{ .GoType}}_{{ .MethodName }}\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return fmt.Errorf(\"%s: %v\", \"{{ .GoType}}_{{ .MethodName }}\",err)")
	fmt.Fprintln(fgo, "    }")
{{ end }}
{{ range .SafeMethods }}
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
	fmt.Fprintln(fgo, "    var crc uint64")
	fmt.Fprintln(fgo, "    _, _, _ = syscall.SyscallN(getcrc,uintptr(unsafe.Pointer(&crc)))")	
    fmt.Fprintln(fgo, "    if crc != {{ .CRC }} {")
    fmt.Fprintln(fgo, "        return fmt.Errorf(\"CRC mismatch %s != %x. DLL is not from same build than go code.\",\"{{ .CRC }}\", crc)")
    fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate__getError, err = syscall.GetProcAddress(dll, \"GetError\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "	   return nil")
	fmt.Fprintln(fgo, "}")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func {{ $.ModuleName }}_getError(rc uintptr)(err error) {")
	fmt.Fprintln(fgo, "    errText := make([]byte, 0, 512)")
	fmt.Fprintln(fgo, "    _, _, _ = syscall.SyscallN(_{{ $.ModuleName }}_gate__getError, rc, uintptr(unsafe.Pointer(&errText)))")
	fmt.Fprintln(fgo, "    return errors.New(string(errText))")
	fmt.Fprintln(fgo, "}")
	fmt.Fprintln(fgo, "")
{{ range .Methods }}
	fmt.Fprintln(fgo, "func (r *{{ .GoType}}) {{ .MethodName }}() (err error) {")
	fmt.Fprint(fgo, "    rc, _, _ := syscall.SyscallN(_{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, uintptr(unsafe.Pointer(r)), ")
    fmt.Fprintln(fgo, "      uintptr(", reflect.TypeOf(new({{ .GoType}})).Elem().Size(), "))")
    fmt.Fprintln(fgo, "    if rc != 0 {")
    fmt.Fprintln(fgo, "         return {{ $.ModuleName }}_getError(rc)")
    fmt.Fprintln(fgo, "    }")
    fmt.Fprintln(fgo, "    return nil")
	fmt.Fprintln(fgo, "}")
    fmt.Fprintln(fgo, "")
{{ end }}
{{ range .SafeMethods }}
	fmt.Fprintln(fgo, "func (r *{{ .GoType}}) {{ .MethodName }}() (err error) {")
	fmt.Fprintln(fgo, "    rc := _{{ $.ModuleName}}_fastcall(_{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, uintptr(unsafe.Pointer(r)), ")
    fmt.Fprintln(fgo, "      uintptr(", reflect.TypeOf(new({{ .GoType}})).Elem().Size(), "))")
    fmt.Fprintln(fgo, "    if rc != 0 {")
    fmt.Fprintln(fgo, "         return {{ $.ModuleName }}_getError(rc)")
    fmt.Fprintln(fgo, "    }")
    fmt.Fprintln(fgo, "    return nil")
	fmt.Fprintln(fgo, "}")
    fmt.Fprintln(fgo, "")
{{ end }}

}
`

const winFastcall = `
func __winFastcall() {
	fasm, err := os.Create("{{ .GoTargetFile }}_windows_amd64.s")
	if err != nil {
		log.Fatal(err)
	}
	defer fasm.Close()

	fmt.Fprintln(fasm)
	fmt.Fprintln(fasm, "TEXT ·_{{ .ModuleName }}_fastcall(SB), 0, $65536-32")
    fmt.Fprintln(fasm, "    MOVQ ptr+8(FP), CX")
    fmt.Fprintln(fasm, "    MOVQ size+16(FP), DX")
    fmt.Fprintln(fasm, "    MOVQ trap+0(FP), AX")
    fmt.Fprintln(fasm, "    MOVQ SP, BX")
    fmt.Fprintln(fasm, "    ADDQ $65472, SP // SP - 4 * 8")
    fmt.Fprintln(fasm, "    ANDQ $~15, SP")
    fmt.Fprintln(fasm, "    CALL AX")
    fmt.Fprintln(fasm, "    MOVQ BX, SP")
    fmt.Fprintln(fasm, "    MOVQ AX, ret+24(FP)")
    fmt.Fprintln(fasm, "    RET")
	fmt.Fprintln(fasm, "")
}
`

const linuxFastcall = `
func __linuxFastcall() {
	fasm, err := os.Create("{{ .GoTargetFile }}_linux_amd64.s")
	if err != nil {
		log.Fatal(err)
	}
	defer fasm.Close()

	fmt.Fprintln(fasm)
	fmt.Fprintln(fasm, "TEXT ·_{{ .ModuleName }}_fastcall(SB), 0, $65536-32")
    fmt.Fprintln(fasm, "    MOVQ ptr+8(FP), DI")
    fmt.Fprintln(fasm, "    MOVQ size+16(FP), SI")
    fmt.Fprintln(fasm, "    MOVQ trap+0(FP), AX")
 	fmt.Fprintln(fasm, "    ADDQ $65504, SP // SP - 4 * 8")
    fmt.Fprintln(fasm, "    CALL AX")
    fmt.Fprintln(fasm, "    ADDQ $-65504, SP")
    fmt.Fprintln(fasm, "    MOVQ AX, ret+24(FP)")
    fmt.Fprintln(fasm, "    RET")
	fmt.Fprintln(fasm, "")
}

`

const loadLinux = `
func __linuxLoader() {
	fgo, err := os.Create("{{ .GoTargetFile }}_linux_amd64.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fgo.Close()
	
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "//")
	fmt.Fprintln(fgo, "package {{ .PackageName }}")
	fmt.Fprintln(fgo, "// Generated file. Not not edit")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "import \"github.com/lakal3/dllcall/linux/syscall\"")
	fmt.Fprintln(fgo, "import \"unsafe\"")
	fmt.Fprintln(fgo, "import \"errors\"")
	fmt.Fprintln(fgo, "import \"fmt\"")
	{{ if .Pin }}
	fmt.Fprintln(fgo, "import \"runtime\"")
	{{ end }}
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "var (")
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate__getError uintptr")
{{ range .Methods }}
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }} uintptr")
{{ end }}
{{ range .SafeMethods }}
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }} uintptr")
{{ end }}
	fmt.Fprintln(fgo, ")")
{{ if .SafeMethods }}
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func _{{ .ModuleName }}_fastcall(trap uintptr, ptr uintptr, size uintptr)(errPtr uintptr)")
	fmt.Fprintln(fgo, "func _{{ .ModuleName }}_fc_alloc()(ret uintptr)")
{{ end }}
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func load_{{ .ModuleName }}(dllPath string)(err error) {")
	fmt.Fprintln(fgo, "    if _{{ $.ModuleName }}_gate__getError != 0 {")
	fmt.Fprintln(fgo, "        return nil")
	fmt.Fprintln(fgo, "    }")
	{{ if not .Pin }}
	fmt.Fprintln(fgo, "    err = syscall.DisableCgocheck()")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
	{{ end }}
	fmt.Fprintln(fgo, "    dll, err := syscall.LoadLibrary(dllPath)")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
{{ range .Methods }}
	fmt.Fprintln(fgo, "     _{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, err = syscall.GetProcAddress(dll, \"{{ .GoType}}_{{ .MethodName }}\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return fmt.Errorf(\"%s: %v\", \"{{ .GoType}}_{{ .MethodName }}\",err)")
	fmt.Fprintln(fgo, "    }")
{{ end }}
{{ range .SafeMethods }}
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
	fmt.Fprintln(fgo, "    var crc uint64")
	fmt.Fprintln(fgo, "    syscall.SyscallN(getcrc,uintptr(unsafe.Pointer(&crc)))")
    fmt.Fprintln(fgo, "    if crc != {{ .CRC }} {")
    fmt.Fprintln(fgo, "        return fmt.Errorf(\"CRC mismatch %s != %x. DLL is not from same build than go code.\",\"{{ .CRC }}\", crc)")
    fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "    _{{ $.ModuleName }}_gate__getError, err = syscall.GetProcAddress(dll, \"GetError\")")
	fmt.Fprintln(fgo, "	   if err != nil {")
	fmt.Fprintln(fgo, "        return err")
	fmt.Fprintln(fgo, "    }")
	fmt.Fprintln(fgo, "	   return nil")
	fmt.Fprintln(fgo, "}")
	fmt.Fprintln(fgo, "")
	fmt.Fprintln(fgo, "func {{ $.ModuleName }}_getError(rc uintptr)(error) {")
	fmt.Fprintln(fgo, "    errText := make([]byte, 0, 512)")
	fmt.Fprintln(fgo, "    syscall.SyscallN(_{{ $.ModuleName }}_gate__getError, rc, uintptr(unsafe.Pointer(&errText)))")
	fmt.Fprintln(fgo, "    return errors.New(string(errText))")
	fmt.Fprintln(fgo, "}")
	fmt.Fprintln(fgo, "")
	{{ $pin := .Pin }}
	var tt reflect.Type
{{ range .Methods }}
	fmt.Fprintln(fgo, "func (r *{{ .GoType}}) {{ .MethodName }}() (err error) {")
	tt = reflect.TypeOf(new({{ .GoType}}))
	{{ if $pin }}
	_emit_pin(fgo, tt)
	{{ end }}
	fmt.Fprint(fgo, "    rc, _, _ := syscall.SyscallN(_{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, uintptr(unsafe.Pointer(r)), ")
    fmt.Fprintln(fgo, "      uintptr(", tt.Elem().Size(), "))")
    fmt.Fprintln(fgo, "    if rc != 0 {")
    fmt.Fprintln(fgo, "         return {{ $.ModuleName }}_getError(rc)")
    fmt.Fprintln(fgo, "    }")
    fmt.Fprintln(fgo, "    return nil")
	fmt.Fprintln(fgo, "}")
    fmt.Fprintln(fgo, "")
{{ end }}
{{ range .SafeMethods }}
	fmt.Fprintln(fgo, "func (r *{{ .GoType}}) {{ .MethodName }}() (err error) {")
	fmt.Fprintln(fgo, "    rc := _{{ $.ModuleName}}_fastcall(_{{ $.ModuleName }}_gate_{{ .GoType}}_{{ .MethodName }}, uintptr(unsafe.Pointer(r)), ")
    fmt.Fprintln(fgo, "      uintptr(", reflect.TypeOf(new({{ .GoType}})).Elem().Size(), "))")
    fmt.Fprintln(fgo, "    if rc != 0 {")
    fmt.Fprintln(fgo, "         return {{ $.ModuleName }}_getError(rc)")
    fmt.Fprintln(fgo, "    }")
    fmt.Fprintln(fgo, "    return nil")
	fmt.Fprintln(fgo, "}")
    fmt.Fprintln(fgo, "")
{{ end }}

}
`

const generatorCode2 = `

func _gen_CType(typeName string, t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		tn := _gen_CType("", t.Elem())
		return  tn + "* " + typeName
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
		_gen_alias[t.Name()] = t.Name() // Prevent recursive parsing of struct
		sb := &strings.Builder{}
		sb.WriteString(" struct ")
	    sb.WriteString(typeName)
		sb.WriteString(" {\n")
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
{{ range .SafeMethods }}
		if typeName == "{{ .GoType }}" {
			sb.WriteString("    GoError *{{.MethodName }}();\n")
		}
{{ end }}
		sb.WriteString("}")
		sb.WriteString(" ")
		sb.WriteString(typeName)
		return sb.String()
	}
	if t.Kind() == reflect.Int64 {
		return "int64_t " + typeName
	}
	if t.Kind() == reflect.Uint64 {
		return "uint64_t " + typeName
	}
	if t.Kind() == reflect.Int32 {
		return "int32_t " + typeName
	}
	if t.Kind() == reflect.Uint32 {
		return "uint32_t " + typeName
	}
	if t.Kind() == reflect.Int16 {
		return "int16_t " + typeName
	}
	if t.Kind() == reflect.Uint16 {
		return "uint16_t " + typeName
	}
	if t.Kind() == reflect.Int8 {
		return "int8_t " + typeName
	}
	if t.Kind() == reflect.Uint8 {
		return "uint8_t " + typeName
	}
	if t.Kind() == reflect.Float64 {
		return "double " + typeName
	}
	if t.Kind() == reflect.Float32 {
		return "float " + typeName
	}

	log.Fatal("Unsupported type ", t.Kind())
	return t.Name()
}

func _emit_pin(fgo *os.File, t reflect.Type) {
	needPin := false
	t = t.Elem()
	for idx := 0; idx < t.NumField(); idx++ {
		f := t.Field(idx)
		if f.Type.Kind() == reflect.Slice || f.Type.Kind() == reflect.Ptr || f.Type.Kind() == reflect.String {
			if !needPin {
				_, _ = fmt.Fprintln(fgo, "    var p runtime.Pinner")
				_, _ = fmt.Fprintln(fgo, "    defer p.Unpin()")
				needPin = true
			}
			switch (f.Type.Kind()) {
			case reflect.Ptr:
				_, _ = fmt.Fprintf(fgo, "    p.Pin(r.%s)\n", f.Name)
			case reflect.Slice:
				_, _ = fmt.Fprintf(fgo, "    p.Pin(unsafe.SliceData(r.%s))\n", f.Name)
			case reflect.String:
				_, _ = fmt.Fprintf(fgo, "    p.Pin(unsafe.StringData(r.%s))\n", f.Name)
			}
		}
	}
}

var _gen_alias = map[string]string {
{{ range .Aliases}}    "{{ .GoType }}" : "{{ .CAlias }}",
{{ end }}
}
`

const header = `
// Generated file. Not not edit
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

`
