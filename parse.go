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
	typeName    string
	comment     string
	methods     []string
	safeMethods []string
	// imports  []string
	ctype string
	cdecl string
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

var methods = []string{"#cmethod", "#csafe_method", "#ctype", "#c"}

func (info *typeInfo) parseComment() error {
	var sb *strings.Builder

	for _, line := range strings.Split(info.comment, "\n") {
		idx := 0
		pos := -1
		for ; idx < len(methods); idx++ {
			pos = strings.Index(line, methods[idx])
			if pos >= 0 {
				break
			}
		}

		if idx == 3 {
			return fmt.Errorf("Invalid #c directive at line: %s", line)
		}
		if sb != nil {
			if pos >= 0 {
				return errors.New("#ctype or #cmethod inside #ctype declaration")
			}
			idx1 := strings.Index(line, "*/")
			if idx1 < 0 {
				idx1 = len(line)
			}
			sb.WriteString(line[:idx1])
			sb.WriteRune('\n')
			continue
		}
		if idx == 0 { // #cmethod
			for _, m := range strings.Split(line[pos+8:], ",") {
				info.methods = append(info.methods, strings.Trim(m, " \t"))
			}
			continue
		}
		if idx == 1 { // #csafe_method
			for _, m := range strings.Split(line[pos+13:], ",") {
				info.safeMethods = append(info.safeMethods, strings.Trim(m, " \t"))
			}
			continue
		}
		if idx == 2 { // #ctype
			info.ctype = strings.Trim(line[pos+6:], " \t")
			sb = &strings.Builder{}
		}

	}
	if sb != nil {
		info.cdecl = sb.String()
	}
	return nil
}
