package main

import (
	"bytes"
	"fmt"
)

type FieldName string

type FieldType byte

const (
	OneField FieldType = iota
	AllFields
	StrictNewField
	MaybeNewField
)

type ExtFieldName struct {
	Name FieldName
	Type FieldType
}

type VarName []FieldName

type ExtVarName []ExtFieldName

type FieldInfo struct{
	varIdx int
	elemIdx int
	elemWasLast bool
	paths Paths
}

func (fieldInfo FieldInfo) print(tab string) {
	fmt.Printf("%svarIdx = %d; elemIdx = %d; elemWasLast = %#v; paths = {\n", tab, fieldInfo.varIdx, fieldInfo.elemIdx, fieldInfo.elemWasLast)
	for fn, fi := range fieldInfo.paths {
		fmt.Printf("%s %s:\n", tab, fn)
		fi.print(tab + "  ")
	}
	fmt.Println(tab + "}")
}

type Paths map[FieldName]*FieldInfo

type GetStatus byte

const (
	IsNotLast GetStatus = iota
	IsLast
	WithError
)

type SetType byte

const (
	Define SetType = iota
	Restore
	Append
)

type SetStatus bool

type Value interface {
	fmt.Stringer
	getNextString(VarName, int, *FieldInfo) (String, GetStatus)
	setString(ExtVarName, String, SetType) (Value, SetStatus)
}
// `Array`, `String` or `Struct`

type Array []Value

func (array Array) String() string {
	if array == nil {
		return "Array"
	}
	return fmt.Sprintf("Array(length = %d)", len(array))
}

func (array Array) getNextString(varName VarName, varIdx int, fieldInfo *FieldInfo) (String, GetStatus) {
	if fieldInfo.varIdx != varIdx {
		return array[fieldInfo.elemIdx].getNextString(varName, varIdx, fieldInfo)
	}
	if fieldInfo.elemWasLast {
		fieldInfo.elemIdx++
		fieldInfo.elemWasLast = false
		if fieldInfo.elemIdx == len(array) {
			fieldInfo.elemIdx = 0
			return nil, WithError
		}
	}
	result, status := array[fieldInfo.elemIdx].getNextString(varName, varIdx, fieldInfo)
	fieldInfo.elemWasLast = status != IsNotLast
	if status != WithError {
		return result, IsNotLast
	}
	return array.getNextString(varName, varIdx, fieldInfo)
}

func (array Array) setString(extVarName ExtVarName, value String, setType SetType) (Value, SetStatus) {
	if len(extVarName) == 0 {
		switch setType {
		case Define:
			return nil, false
		case Restore:
			return value, true
		case Append:
			return append(array, value), true
		default:
			return nil, false
		}
	}
	switch extVarName[0].Type {
	case StrictNewField:
		if setType == Define {
			return nil, false
		}
		newStruct := make(Struct)
		_, _ = newStruct.setString(extVarName, value, Define)
		switch setType {
		case Restore:
			return newStruct, true
		case Append:
			return append(array, newStruct), true
		default:
			return nil, false
		}
	case MaybeNewField:
		extVarName[0].Type = OneField
		newValue, ok := array.setString(extVarName, value, Define)
		if ok {
			return newValue, true
		}
		extVarName[0].Type = StrictNewField
		return array.setString(extVarName, value, setType)
	case OneField:
		i := len(array) - 1
		newValue, ok := array[i].setString(extVarName, value, setType)
		if newValue != nil {
			array[i] = newValue
		}
		return nil, ok
	case AllFields:
		for i := range array {
			newValue, _ := array[i].setString(extVarName, value, setType)
			if newValue != nil {
				array[i] = newValue
			}
		}
		return nil, true
	default:
		return nil, false
	}
}

type String []byte

func (string_ String) String() string {
	if string_ == nil {
		return "String"
	}
	return fmt.Sprintf("String(%s)", string(string_))
}

func (string_ String) getNextString(varName VarName, _ int, _ *FieldInfo) (String, GetStatus) {
	if len(varName) == 0 {
		return string_, IsLast
	}
	return nil, WithError
}

func (string_ String) setString(extVarName ExtVarName, value String, setType SetType) (Value, SetStatus) {
	switch setType {
	case Define:
		return nil, false
	case Restore:
		if len(extVarName) == 0 {
			return value, true
		}
		newStruct := make(Struct)
		_, _ = newStruct.setString(extVarName, value, Define)
		return newStruct, true
	case Append:
		if len(extVarName) == 0 {
			return Array{string_, value}, true
		}
		newStruct := make(Struct)
		_, _ = newStruct.setString(extVarName, value, Define)
		return Array{string_, newStruct}, true
	default:
		return nil, false
	}
}

func (string_ String) print(tab string) {
	fmt.Println(tab + string_.String())
}

func (string_ String) JoinStrings(strings []String) String {
	if len(strings) == 0 {
		return nil
	}
	var buf bytes.Buffer
	buf.Write(strings[0])
	for i := 1; i < len(strings); i++ {
		buf.Write(string_)
		buf.Write(strings[i])
	}
	return buf.Bytes()
}

type Struct map[FieldName]Value

func (struct_ Struct) String() string {
	return "Struct"
}

func (struct_ Struct) getNextString(varName VarName, varIdx int, fieldInfo *FieldInfo) (String, GetStatus) {
	if len(varName) == 0 {
		return nil, WithError
	}
	fieldName := varName[0]
	restVarName := varName[1:]
	if _, ok := struct_[fieldName]; !ok {
		return nil, WithError
	}
	if fieldInfo.paths == nil {
		fieldInfo.paths = make(Paths)
	}
	if _, ok := fieldInfo.paths[fieldName]; !ok {
		fieldInfo.paths[fieldName] = &FieldInfo{varIdx: varIdx}
	}
	return struct_[fieldName].getNextString(restVarName, varIdx, fieldInfo.paths[fieldName])
}

func (struct_ Struct) setString(extVarName ExtVarName, value String, setType SetType) (Value, SetStatus) {
	if len(extVarName) == 0 {
		switch setType {
		case Define:
			return nil, false
		case Restore:
			return value, true
		case Append:
			return Array{struct_, value}, true
		default:
			return nil, false
		}
	}
	fieldName := extVarName[0].Name
	restExtVarName := extVarName[1:]
	if _, ok := struct_[fieldName]; !ok {
		if len(restExtVarName) == 0 {
			struct_[fieldName] = value
		} else {
			struct_[fieldName] = make(Struct)
			_, _ = struct_[fieldName].setString(restExtVarName, value, setType)
		}
		return nil, true
	}
	newValue, ok := struct_[fieldName].setString(restExtVarName, value, setType)
	if !ok {
		return nil, false
	}
	if newValue != nil {
		struct_[fieldName] = newValue
	}
	return nil, true
}

func (struct_ Struct) getStrings(varNames []VarName, varIdx int, fieldInfo *FieldInfo, result *[][]String) int {
	if varIdx == len(varNames) {
		return 1
	}
	status := IsNotLast
	var string_ String
	recordsCount := 0
	for status == IsNotLast {
		string_, status = struct_.getNextString(varNames[varIdx], varIdx, fieldInfo)
		if status == WithError {
			break
		}
		subrecordsCount := struct_.getStrings(varNames, varIdx+1, fieldInfo, result)
		recordsCount += subrecordsCount
		for i := 0; i < subrecordsCount; i++ {
			(*result)[varIdx] = append((*result)[varIdx], string_)
		}
	}
	return recordsCount
}

func (struct_ Struct) GetStrings(varNames []VarName) [][]String {
	result := make([][]String, len(varNames))
	_ = struct_.getStrings(varNames, 0, &FieldInfo{varIdx: -1}, &result)
	return result
}

func (struct_ Struct) SetStrings(extVarName ExtVarName, values []String) {
	_, _ = struct_.setString(extVarName, String(nil).JoinStrings(values), Restore)
}

func (struct_ Struct) AddStrings(extVarName ExtVarName, values []String) {
	_, _ = struct_.setString(extVarName, String(nil).JoinStrings(values), Append)
}

func ToGoStruct(value Value) interface{} {
	switch value := value.(type) {
	case Struct:
		result := make(map[string]interface{})
		for field, subvalue := range value {
			result[string(field)] = ToGoStruct(subvalue)
		}
		return result
	case Array:
		result := make([]interface{}, len(value))
		for i, subvalue := range value {
			result[i] = ToGoStruct(subvalue)
		}
		return result
	case String:
		return string(value)
	default:
		return nil
	}
}

func FromGoStruct(goValue interface{}) Value {
	switch goValue := goValue.(type) {
	case map[string]interface{}:
		result := make(Struct)
		for key, value := range goValue {
			result[FieldName(key)] = FromGoStruct(value)
		}
		return result
	case []interface{}:
		result := make(Array, len(goValue))
		for i, value := range goValue {
			result[i] = FromGoStruct(value)
		}
		return result
	case string:
		return String(goValue)
	default:
		return nil
	}
}
