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
	GoTypeByte
	GoTypeRune
	GoTypeString
	GoTypeBool
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
	case GoTypeByte:   return "byte"
	case GoTypeRune:   return "rune"
	case GoTypeString: return "string"
	case GoTypeBool:   return "bool"
	default: return ""
	}
}

func (t GoBasicType) StringAsKleeType() string {
	switch t {
	case GoTypeInt:    fallthrough
	case GoTypeInt8:   fallthrough
	case GoTypeInt16:  fallthrough
	case GoTypeInt32:  fallthrough
	case GoTypeInt64:  fallthrough
	case GoTypeRune:   return "int"
	case GoTypeUInt:   fallthrough
	case GoTypeUInt8:  fallthrough
	case GoTypeUInt16: fallthrough
	case GoTypeUInt32: fallthrough
	case GoTypeUInt64: fallthrough
	case GoTypeByte:   return "uint"
	case GoTypeString: return "text"
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
	"byte":   GoTypeByte,
	"rune":   GoTypeRune,
	"string": GoTypeString,
	"bool":   GoTypeBool,
}

type GoVar struct {
	Name string
	Type GoType
}

func TypeAsString(t ast.Expr) string {
	if tIdent, ok := t.(*ast.Ident); ok {
		return tIdent.String()
	} else if tArray, ok := t.(*ast.ArrayType); ok {
		return fmt.Sprintf("[]%s", TypeAsString(tArray.Elt))
	} else if tMap, ok := t.(*ast.MapType); ok {
		return fmt.Sprintf("map[%s]%s", TypeAsString(tMap.Key), TypeAsString(tMap.Value))
	} else {
		return ""
	}
}

func (function GoFunction) Params() (res []GoVar, err error) {
	return Vars(function.AST.Type.Params.List)
}

func (function GoFunction) Results() (res []GoVar, err error) {
	return Vars(function.AST.Type.Results.List)
}

func Vars(vars []*ast.Field) (res []GoVar, err error) {
	for _, var_ := range vars {
		varTypeIdent, ok := var_.Type.(*ast.Ident)
		if ! ok {
			err = fmt.Errorf("unknown type")
			return
		}
		varType, ok := GoTypes[varTypeIdent.Name]
		if ! ok {
			err = fmt.Errorf("unknown type")
			return
		}
		varNames := var_.Names
		if len(varNames) == 0 {
			varNames = []*ast.Ident{{Name: ""}}
		}
		for _, varName := range varNames {
			res = append(res, GoVar{varName.Name, varType})
		}
	}
	return
}
