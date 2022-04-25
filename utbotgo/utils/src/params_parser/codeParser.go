package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
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

type GoType interface {
	StringAsGoType() string
	StringAsKleeType() string
}

type GoBasicType int

const (
	GoTypeInt GoBasicType = iota
	GoTypeInt8
	GoTypeInt16
	GoTypeInt32
	GoTypeInt64
	GoTypeUInt
	GoTypeUInt8
	GoTypeUInt16
	GoTypeUInt32
	GoTypeUInt64
)

func (t GoBasicType) StringAsGoType() string {
	switch t {
	case GoTypeInt:    return "int"
	case GoTypeInt8:   return "int8"
	case GoTypeInt16:  return "int16"
	case GoTypeInt32:  return "int32"
	case GoTypeInt64:  return "int64"
	case GoTypeUInt:   return "uint"
	case GoTypeUInt8:  return "uint8"
	case GoTypeUInt16: return "uint16"
	case GoTypeUInt32: return "uint32"
	case GoTypeUInt64: return "uint64"
	default: return ""
	}
}

func (t GoBasicType) StringAsKleeType() string {
	switch t {
	case GoTypeInt:    fallthrough
	case GoTypeInt8:   fallthrough
	case GoTypeInt16:  fallthrough
	case GoTypeInt32:  fallthrough
	case GoTypeInt64:  return "int"
	case GoTypeUInt:   fallthrough
	case GoTypeUInt8:  fallthrough
	case GoTypeUInt16: fallthrough
	case GoTypeUInt32: fallthrough
	case GoTypeUInt64: return "uint"
	default: return ""
	}
}

type GoStructType map[string]GoType

func (t GoStructType) StringAsGoType() string {
	var fields []string
	for fieldName, fieldType := range t {
		fields = append(fields, fmt.Sprintf("%s %s", fieldName, fieldType.StringAsGoType()))
	}
	return fmt.Sprintf("struct{%s}", strings.Join(fields, ";"))
}

func (_ GoStructType) StringAsKleeType() string {
	return "b64"
}

var GoTypes = map[string]GoType{
	"int":    GoTypeInt,
	"int8":   GoTypeInt8,
	"int16":  GoTypeInt16,
	"int32":  GoTypeInt32,
	"int64":  GoTypeInt64,
	"uint":   GoTypeUInt,
	"uint8":  GoTypeUInt8,
	"uint16": GoTypeUInt16,
	"uint32": GoTypeUInt32,
	"uint64": GoTypeUInt64,
}

type GoParam struct {
	Name string
	Type GoType
}

func (function GoFunction) Params() (res []GoParam, err error) {
	params := function.AST.Type.Params.List
	for _, param := range params {
		paramTypeIdent, ok := param.Type.(*ast.Ident)
		if ! ok {
			err = fmt.Errorf("unknown type")
			return
		}
		paramType, ok := GoTypes[paramTypeIdent.Name]
		if ! ok {
			err = fmt.Errorf("unknown type")
			return
		}
		for _, paramName := range param.Names {
			res = append(res, GoParam{paramName.Name, paramType})
		}
	}
	return
}

func (function GoFunction) ResultTypeAsString() string {
	results := function.AST.Type.Results.List
	if len(results) == 1 {
		return results[0].Type.(fmt.Stringer).String()
	} else {
		resultTypes := make([]string, len(results))
		for i := range results {
			resultTypes[i] = results[i].Type.(fmt.Stringer).String()
		}
		return fmt.Sprintf("(%s)", strings.Join(resultTypes, ", "))
	}
}
