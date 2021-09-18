package main

import (
    "go/parser"
    "go/token"
    "go/ast"
)

func NewFile (fileName string) (File, error) {
    fset := token.NewFileSet()
    file, err := parser.ParseFile(fset, fileName, nil, 0)
    if err != nil {
        return File{nil, fileName}, err
    }
    return File{file, fileName}, nil
}

func (file File) PackageName () string {
    return file.AST.Name.Name
}

func (file File) Functions () []Function {
    res := []Function{}
    for _, decl := range file.AST.Decls {
        if function, ok := decl.(*ast.FuncDecl); ok {
            res = append(res, Function{function})
        }
    }
    return res
}

func (function Function) Name () string {
    return function.AST.Name.Name
}

func (function Function) ParamTypes () (res []Type) {
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
