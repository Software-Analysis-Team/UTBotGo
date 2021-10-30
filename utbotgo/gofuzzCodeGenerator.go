package main

import (
	"fmt"
	fuzz "github.com/google/gofuzz"
	"strings"
)

func GofuzzTestingCode(file GoFile) string {
	tb := &TextBuilder{}
	GofuzzPackageDeclaration(tb, file)
	GofuzzImportedLibraries(tb)
	GofuzzFuzzingForFile(tb, file, 10)
	return tb.String()
}

func GofuzzPackageDeclaration(tb *TextBuilder, file GoFile) {
	tb.Section(func() {
		tb.WriteLine("package " + file.PackageName())
	})
}

func GofuzzImportedLibraries(tb *TextBuilder) {
	tb.Section(func() {
		tb.Subsection("import (", ")", func() {
			tb.WriteLine("\"fmt\"")
			tb.WriteLine("\"strconv\"")
			tb.WriteLine("\"testing\"")
		})
	})
}

func GofuzzFuzzingForFile(tb *TextBuilder, file GoFile, repetitionsNumber int) {
	for _, function := range file.Functions() {
		GofuzzFuzzingForFunction(tb, function, repetitionsNumber)
	}
}

func GofuzzFuzzingForFunction(tb *TextBuilder, function GoFunction, repetitionsNumber int) {
	tb.Section(func() {
		funcDescription := fmt.Sprintf("func Test_%s(t *testing.T) {", function.Name())
		tb.Subsection(funcDescription, "}", func() {
			dataName := GofuzzInit(tb, function, repetitionsNumber)
			GofuzzMainLoop(tb, function, repetitionsNumber, dataName)
		})
	})
}

func GofuzzInit(tb *TextBuilder, function GoFunction, repetitionsNumber int) string {
	var dataName string
	tb.Section(func() {
		argsType := GofuzzArgsType(tb, function)
		dataName = GofuzzDataDefinition(tb, function, repetitionsNumber, argsType)
	})
	return dataName
}

func GofuzzArgsType(tb *TextBuilder, function GoFunction) string {
	argsTypeName := "args"
	tb.Section(func() {
		tb.Subsection(fmt.Sprintf("type %s struct {", argsTypeName), "}", func() {
			for _, param := range function.Params() {
				tb.WriteLine(param.Name + " " + param.Type)
			}
		})
	})
	return argsTypeName
}

type Variable struct {
	Value interface{}
}

func (v *Variable) New(type_ string) {
	switch type_ {
	case "int":
		v.Value = new(int)
	case "string":
		v.Value = new(string)
	}
}

func (v *Variable) String() (res string) {
	switch val := v.Value.(type) {
	case *int:
		res = fmt.Sprintf("%d", *val)
	case *string:
		res = fmt.Sprintf("%#v", *val)
	}
	return
}

func GofuzzDataDefinition(tb *TextBuilder, function GoFunction, repetitionsNumber int, argsType string) string {
	fuzzer := fuzz.New()
	variables := make([]Variable, len(function.Params()))
	funcParams := function.Params()
	dataName := "data"
	dataDefinition := fmt.Sprintf("%s := []%s{", dataName, argsType)
	tb.Section(func() {
		tb.Subsection(dataDefinition, "}", func() {
			for j := 0; j < repetitionsNumber; j++ {
				tb.Subsection("{", "},", func() {
					for i := range funcParams {
						variables[i].New(funcParams[i].Type)
						fuzzer.Fuzz(variables[i].Value)
						tb.WriteLine(funcParams[i].Name + ": " + variables[i].String() + ",")
					}
				})
			}
		})
	})
	return dataName
}

func GofuzzMainLoop(tb *TextBuilder, function GoFunction, repetitionsNumber int, dataName string) {
	counterName := "i"
	tb.Section(func() {
		forDefinition := fmt.Sprintf(
			"for %s := 0; %s < %d; %s++ {",
			counterName,
			counterName,
			repetitionsNumber,
			counterName,
		)
		tb.Subsection(forDefinition, "}", func() {
			GofuzzMainLoopBlock(tb, function, dataName, counterName)
		})
	})
}

func GofuzzMainLoopBlock(tb *TextBuilder, function GoFunction, dataName string, counterName string) {
	paramNames := make([]string, len(function.Params()))
	for i, param := range function.Params() {
		paramNames[i] = fmt.Sprintf("%s[%s].%s", dataName, counterName, param.Name)
	}
	tb.Section(func() {
		subtestRunning := fmt.Sprintf("t.Run(strconv.Itoa(%s), func(t *testing.T) {", counterName)
		tb.Subsection(subtestRunning, "})", func() {
			tb.Section(func() {
				tb.Subsection("defer func() {", "}()", func() {
						tb.Subsection("if recover() != nil {", "}", func() {
						errorMessage := GofuzzErrorMessage(function, paramNames)
						tb.WriteLine(fmt.Sprintf("t.Error(%s)", errorMessage))
					})
				})
			})
			GofuzzCallingFunction(tb, function, paramNames)
		})
	})
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

func GofuzzCallingFunction(tb *TextBuilder, function GoFunction, paramNames []string) {
	tb.Section(func() {
		tb.WriteLine(function.Name() + "(" + strings.Join(paramNames, ", ") + ")")
	})
}
