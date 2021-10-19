package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

type GoFile struct {
	AST  *ast.File
	Name string
}

func ParseGoFile(fileName string) (GoFile, error) {
	fileSet := token.NewFileSet()
	file, err := parser.ParseFile(fileSet, fileName, nil, 0)
	if err != nil {
		return GoFile{nil, fileName}, err
	}
	return GoFile{file, fileName}, nil
}

func (file GoFile) PackageName() string {
	return file.AST.Name.Name
}

type GoFunction struct {
	AST *ast.FuncDecl
}

func (file GoFile) Functions() (res []GoFunction) {
	for _, decl := range file.AST.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			res = append(res, GoFunction{function})
		}
	}
	return res
}

func (function GoFunction) Name() string {
	return function.AST.Name.Name
}

type Param struct {
	Name string
	Type string
}

func (function GoFunction) Params() (res []Param) {
	params := function.AST.Type.Params.List
	for _, param := range params {
		paramType := param.Type.(*ast.Ident).Name
		for _, paramName := range param.Names {
			res = append(res, Param{paramName.Name, paramType})
		}
	}
	return
}
