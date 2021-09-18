package main

import "go/ast"

type File struct {
    AST *ast.File
    Name string
}

type Function struct {
    AST *ast.FuncDecl
}

type Type string

type Data string

type Test []Data

type TestSet map[string][]Test
