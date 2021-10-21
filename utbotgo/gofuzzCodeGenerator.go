package main

import (
	"fmt"
	"strings"
)

func GofuzzTestingCode(file GoFile) string {
	tb := &TextBuilder{}
	GofuzzPackageDeclaration(tb, file)
	GofuzzImportedLibraries(tb)
	GofuzzFuzzingForFile(tb, file)
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
			tb.WriteLine("fuzz \"github.com/google/gofuzz\"")
			tb.WriteLine("\"strconv\"")
			tb.WriteLine("\"testing\"")
		})
	})
}

func GofuzzFuzzingForFile(tb *TextBuilder, file GoFile) {
	for _, function := range file.Functions() {
		GofuzzFuzzingForFunction(tb, function)
	}
}

func GofuzzFuzzingForFunction(tb *TextBuilder, function GoFunction) {
	tb.Section(func() {
		funcDescription := fmt.Sprintf("func Test_%s(t *testing.T) {", function.Name())
		tb.Subsection(funcDescription, "}", func() {
			paramNames, fuzzerName := GofuzzInit(tb, function)
			GofuzzMainLoop(tb, function, 10, paramNames, fuzzerName)
		})
	})
}

func GofuzzInit(tb *TextBuilder, function GoFunction) ([]string, string) {
	var paramName string
	var paramNames []string
	fuzzerName := "fuzzer"
	tb.Section(func() {
		tb.Subsection("var (", ")", func() {
			for _, param := range function.Params() {
				paramName = param.Name + "_"
				tb.WriteLine(paramName + " " + param.Type)
				paramNames = append(paramNames, paramName)
			}
		})
		tb.WriteLine(fuzzerName + " := fuzz.New()")
	})
	return paramNames, fuzzerName
}

func GofuzzMainLoop(tb *TextBuilder, function GoFunction, repetitionsNumber int, paramNames []string, fuzzerName string) {
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
			GofuzzMainLoopBlock(tb, function, paramNames, fuzzerName, counterName)
		})
	})
}

func GofuzzMainLoopBlock(tb *TextBuilder, function GoFunction, paramNames []string, fuzzerName string, counterName string) {
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
			GofuzzFuzzingSection(tb, paramNames, fuzzerName)
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

func GofuzzFuzzingSection(tb *TextBuilder, paramNames []string, fuzzerName string) {
	tb.Section(func() {
		for _, paramName := range paramNames {
			tb.WriteLine(fmt.Sprintf("%s.Fuzz(&%s)", fuzzerName, paramName))
		}
	})
}

func GofuzzCallingFunction(tb *TextBuilder, function GoFunction, paramNames []string) {
	tb.Section(func() {
		tb.WriteLine(function.Name() + "(" + strings.Join(paramNames, ", ") + ")")
	})
}
