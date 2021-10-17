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
	AST  *ast.File
	Name string
}

type GoFunction struct {
	AST *ast.FuncDecl
}

type Param struct {
	Name string
	Type Type
}

type Type string

type Data string

type Test []Data

type TestSet map[string][]Test

type FuzzingMode int

const (
	Gofuzz FuzzingMode = iota
	Go_fuzz
	TestingFuzz
)

type GeneratingResult struct {
	GeneratingFilesAmount int
}

type ParsedArgs struct {
	Object FSObject
	Mode   FuzzingMode
}
