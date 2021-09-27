package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func NewGoFile(fileName string) (GoFile, error) {
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

func (file GoFile) Functions() []GoFunction {
	var res []GoFunction
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

func (function GoFunction) ParamTypes() (res []Type) {
	res = []Type{}
	params := function.AST.Type.Params.List
	for _, param := range params {
		paramType := param.Type.(*ast.Ident).Name
		for range param.Names {
			res = append(res, Type(paramType))
		}
	}
	return
}
