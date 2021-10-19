package main

import (
	"fmt"
	"strings"
)

func GofuzzTestingCode(file GoFile) string {
	return GofuzzPackageDeclaration(file) +
		GofuzzImportedLibraries() +
		GofuzzFuzzingForFile(file)
}

func GofuzzPackageDeclaration(file GoFile) string {
	return fmt.Sprintf("package %s\n", file.PackageName())
}

func GofuzzImportedLibraries() string {
	temp := `
import (
	"fmt"
	fuzz "github.com/google/gofuzz"
	"strconv"
	"testing"
)

`
	return temp
}

func GofuzzFuzzingForFile(file GoFile) (code string) {
	for _, function := range file.Functions() {
		code += GofuzzFuzzingForFunction(function)
	}
	return
}

func GofuzzFuzzingForFunction(function GoFunction) string {
	funcDecl := GofuzzFunctionDecl(function)
	paramsDefinition, paramNames, fuzzerName := GofuzzInit(function)
	forLoop := GofuzzMainLoop(function, 10, paramNames, fuzzerName)
	return funcDecl + " {\n" + paramsDefinition + forLoop + "}\n\n"
}

func GofuzzFunctionDecl(function GoFunction) string {
	return fmt.Sprintf("func Test_%s (t *testing.T)", function.Name())
}

func GofuzzInit(function GoFunction) (string, []string, string) {
	var resultBuilder strings.Builder
	var paramNames []string
	fuzzerName := "fuzzer"
	resultBuilder.WriteString("\tvar (\n")
	for _, param := range function.Params() {
		resultBuilder.WriteString(fmt.Sprintf("\t\t%s_ %s\n", param.Name, param.Type))
		paramNames = append(paramNames, fmt.Sprintf("%s_", param.Name))
	}
	resultBuilder.WriteString("\t)\n")
	resultBuilder.WriteString(fmt.Sprintf("\t%s := fuzz.New()\n", fuzzerName))
	return resultBuilder.String(), paramNames, fuzzerName
}

func GofuzzMainLoop(function GoFunction, repetitionsNumber int, paramNames []string, fuzzerName string) string {
	counterName := "i"
	temp := `	for %s := 0; %s < %d; %s++ {
%s
	}
`
	return fmt.Sprintf(
		temp,
		counterName,
		counterName,
		repetitionsNumber,
		counterName,
		GofuzzMainLoopBlock(function, paramNames, fuzzerName, counterName),
	)
}

func GofuzzMainLoopBlock(function GoFunction, paramNames []string, fuzzerName string, counterName string) string {
	temp := `		t.Run(strconv.Itoa(%s), func(t *testing.T) {
			defer func() {
				if recover() != nil {
					t.Error(%s)
				}
			}()
%s
			%s
		})`
	errorMessage := GofuzzErrorMessage(function, paramNames)
	fuzzingSection := GofuzzFuzzingSection(paramNames, fuzzerName)
	callingFunction := GofuzzCallingFunction(function, paramNames)
	return fmt.Sprintf(temp, counterName, errorMessage, fuzzingSection, callingFunction)
}

func GofuzzErrorMessage(function GoFunction, paramNames []string) string {
	var resultBuilder strings.Builder
	resultBuilder.WriteString("fmt.Sprintf(\"input args:\\n")
	for _, param := range function.Params() {
		resultBuilder.WriteString(fmt.Sprintf("\\t%s = %%#v\\n", param.Name))
	}
	resultBuilder.WriteString("\"")
	for _, paramName := range paramNames {
		resultBuilder.WriteString(", " + paramName)
	}
	resultBuilder.WriteString(")")
	return resultBuilder.String()
}

func GofuzzFuzzingSection(paramNames []string, fuzzerName string) string {
	var resultBuilder strings.Builder
	for _, paramName := range paramNames {
		resultBuilder.WriteString(fmt.Sprintf("\t\t\t%s.Fuzz(&%s)\n", fuzzerName, paramName))
	}
	return resultBuilder.String()
}

func GofuzzCallingFunction(function GoFunction, paramNames []string) string {
	return function.Name() + "(" + strings.Join(paramNames, ", ") + ")"
}
