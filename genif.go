package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

func generate(gofile, cfile string) error {
	tmpName := ffSlash(cfile) + "_temp.go"
	dt := &genData{
		GoTargetFile: ffSlash(gofile), TargetFile: ffSlash(cfile), Header: header, PackageName: packageName, CRC: crc,
	}
	tmp := filepath.Base(gofile)
	dt.ModuleName = tmp[0 : len(tmp)-5] // _impl.go removed

	dt.Aliases = make(map[string]genAlias)
	for _, st := range stdAlias {
		dt.Aliases[st.GoType] = st
	}
	for _, it := range typeList {
		if len(it.cdecl) > 0 {
			dt.CDecl = append(dt.CDecl, it.cdecl)
		}
		if len(it.ctype) > 0 {
			dt.Aliases[it.typeName] = genAlias{GoType: it.typeName, CAlias: it.ctype}
		} else {
			dt.GoTypes = append(dt.GoTypes, it.typeName)
		}
		if len(it.methods) > 0 {
			for _, m := range it.methods {
				dt.Methods = append(dt.Methods, genMethod{GoType: it.typeName, MethodName: m})
			}
		}
	}
	if len(dt.GoTypes) == 0 {
		return errors.New("No types defines in interface")
	}

	content, err := genTempFile(dt)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(tmpName, content, 0x600)
	if err != nil {
		return err
	}
	err = invokeGen(tmpName)
	if err != nil {
		return err
	}
	if !fKeep {
		os.Remove(tmpName)
	}
	return nil
}

func ffSlash(path string) string {
	return strings.Replace(path, "\\", "/", -1)
}

var stdAlias = []genAlias{
	genAlias{GoType: "uint8", CAlias: "uint8_t"},
	genAlias{GoType: "int8", CAlias: "int8_t"},
	genAlias{GoType: "uint16", CAlias: "uint16_t"},
	genAlias{GoType: "int16", CAlias: "int16_t"},
	genAlias{GoType: "uint32", CAlias: "uint32_t"},
	genAlias{GoType: "int32", CAlias: "int32_t"},
	genAlias{GoType: "uint64", CAlias: "uint64_t"},
	genAlias{GoType: "int64", CAlias: "int64_t"},
	genAlias{GoType: "float32", CAlias: "float"},
	genAlias{GoType: "float64", CAlias: "double"},
	genAlias{GoType: "string", CAlias: "GoString"},
	genAlias{GoType: "uintptr", CAlias: "void *"},
	genAlias{GoType: "bool", CAlias: "bool"},
}

type genAlias struct {
	GoType string
	CAlias string
}

type genMethod struct {
	MethodName string
	GoType     string
}

type genData struct {
	PackageName  string
	GoTargetFile string
	ModuleName   string
	TargetFile   string
	CDecl        []string
	Aliases      map[string]genAlias
	GoTypes      []string
	Header       string
	Methods      []genMethod
	CRC          string
}

var genFuncs = template.FuncMap{
	"Quote": func(str string) string {
		sb := &strings.Builder{}
		for _, r := range str {
			switch r {
			case '\n':
				sb.WriteString("\\n")
			case '\t':
				sb.WriteString("    ")
			case '\r':
			case '"':
				sb.WriteString("\\\"")
			case '\\':
				sb.WriteString("\\\\")
			default:
				sb.WriteRune(r)
			}
		}
		return sb.String()
	},
}

func genTempFile(dt *genData) (content []byte, err error) {
	bf := &bytes.Buffer{}
	fast.Name.Name = "main"
	addImport("fmt")
	addImport("log")
	addImport("os")
	addImport("reflect")
	addImport("strings")
	removeComments()
	err = printer.Fprint(bf, fset, fast)
	if err != nil {
		return nil, err
	}
	templ := template.New("gen")
	templ.Funcs(genFuncs)
	templ, err = templ.Parse(generatorCode1 + loadWin + loadLinux + generatorCode2)
	if err != nil {
		return nil, err
	}
	err = templ.Execute(bf, dt)
	if err != nil {
		return nil, err
	}
	return bf.Bytes(), nil
}

func addImport(path string) {
	for _, imp := range fast.Imports {
		if imp.Path.Value == path {
			return
		}
	}
	imp := &ast.GenDecl{
		Tok: token.IMPORT,
		Specs: []ast.Spec{&ast.ImportSpec{
			Path: &ast.BasicLit{Value: "\"" + path + "\"", Kind: token.STRING},
		}}}
	fast.Decls = append([]ast.Decl{imp}, fast.Decls...)
}

func removeComments() {
	fast.Doc = nil
	fast.Comments = []*ast.CommentGroup{}
}

func invokeGen(genFile string) error {
	cmd := exec.Command("go", "run", genFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Command failed, output: %s", string(output))
	}
	return nil
}
