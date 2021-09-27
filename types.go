package main

import (
    "go/ast"
    "io/fs"
)

type FSObject struct {
    Name string
    Mode fs.FileMode
}

type GoFile struct {
    AST *ast.File
    Name string
}

type GoFunction struct {
    AST *ast.FuncDecl
}

type Type string

type Data string

type Test []Data

type TestSet map[string][]Test
