package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

var fset *token.FileSet
var fast *ast.File

type typeInfo struct {
	typeName string
	comment  string
	methods  []string
	imports  []string
	ctype    string
	cdecl    string
}

var typeMap = make(map[string]*typeInfo)
var typeList = make([]*typeInfo, 0, 10)
var packageName string
var crc string

func parseGoFile(gofile string) (err error) {
	fset = token.NewFileSet()
	fast, err = parser.ParseFile(fset, gofile, nil, parser.ParseComments)
	if err != nil {
		return
	}

	packageName = fast.Name.Name
	for _, dc := range fast.Decls {
		err = parseDecl(dc)
		if err != nil {
			return err
		}
	}
	for _, info := range typeMap {
		err = info.parseComment()
		if err != nil {
			return err
		}
	}
	bytes, err := ioutil.ReadFile(gofile)
	if err != nil {
		return err
	}
	tmp := fmt.Sprintf("%x", md5.Sum(bytes))
	crc = "0x" + tmp[0:16]

	return
}

func parseDecl(decl ast.Decl) error {

	gdl, ok := decl.(*ast.GenDecl)
	if !ok {
		p := fset.Position(decl.Pos())
		return fmt.Errorf("Go file should contain only types, not %T at %d", p.Line)
	}
	comments := addComments("", gdl.Doc)

	for _, spec := range gdl.Specs {
		err := parseSpec(spec, comments)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseSpec(spec ast.Spec, parentComment string) error {
	_, ok := spec.(*ast.ValueSpec)
	if ok {
		return nil // Skip values and constants
	}
	_, ok = spec.(*ast.ImportSpec)
	if ok {
		return nil // Skip imports
	}
	tsp, ok := spec.(*ast.TypeSpec)
	if !ok {
		p := fset.Position(spec.Pos())
		return fmt.Errorf("Invalid specification %T, should be typespec %d", spec, p.Line)
	}
	comments := addComments(parentComment, tsp.Doc)
	tt := &typeInfo{
		typeName: tsp.Name.String(), comment: comments,
	}
	typeMap[tsp.Name.String()] = tt
	typeList = append(typeList, tt)

	return nil
}

func addComments(current string, group *ast.CommentGroup) string {
	if group == nil || len(group.List) == 0 {
		return current
	}
	sb := &strings.Builder{}
	sb.WriteString(current)
	for _, l := range group.List {
		sb.WriteString(l.Text)
		sb.WriteRune('\n')
	}
	return sb.String()
}

func (info *typeInfo) parseComment() error {
	var sb *strings.Builder

	for _, line := range strings.Split(info.comment, "\n") {
		idx1 := strings.Index(line, "#cmethod")
		idx2 := strings.Index(line, "#ctype")
		idx3 := strings.Index(line, "#cimport")
		idx4 := strings.Index(line, "#c")
		if idx4 >= 0 && idx1 < 0 && idx2 < 0 && idx3 < 0 {
			return fmt.Errorf("Invalid #c directive at line: %s", line)
		}
		if sb != nil {
			if idx1 >= 0 || idx2 >= 0 {
				return errors.New("#ctype or #cmethod inside #ctype declaration")
			}
			idx1 = strings.Index(line, "*/")
			if idx1 < 0 {
				idx1 = len(line)
			}
			sb.WriteString(line[:idx1])
			sb.WriteRune('\n')
			continue
		}
		if idx1 >= 0 {
			for _, m := range strings.Split(line[idx1+8:], ",") {
				info.methods = append(info.methods, strings.Trim(m, " \t"))
			}
			continue
		}
		if idx2 >= 0 {
			info.ctype = strings.Trim(line[idx2+6:], " \t")
			sb = &strings.Builder{}
		}
		if idx3 >= 0 {
			info.imports = append(info.imports, strings.Trim(line[idx3+8:], " \t"))
		}
	}
	if sb != nil {
		info.cdecl = sb.String()
	}
	return nil
}
