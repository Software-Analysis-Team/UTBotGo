package main

import (
	"fmt"
	"strings"
)

func PrintPackageDeclaration(file GoFile) string {
	return fmt.Sprintf("package %s\n\n", file.PackageName())
}

func PrintImportedLibraries() string {
	return "import \"testing\"\n\n"
}

func PrintDataFromTest(test Test) string {
	dataStr := make([]string, len(test))
	for i := range dataStr {
		dataStr[i] = string(test[i])
	}
	return strings.Join(dataStr, ", ")
}

func PrintCallingOfTestedFunction(function GoFunction, test Test) string {
	return fmt.Sprintf("%s(%s)", function.Name(), PrintDataFromTest(test))
}

func PrintOneTest(index int, function GoFunction, test Test) string {
	temp := `	t.Run("%d", func() {
		defer func() {
			if recover() != nil {
				t.Fail()
			}
		}()
		%s
	})
`
	return fmt.Sprintf(
		temp,
		index + 1,
		PrintCallingOfTestedFunction(function, test),
	)
}

func PrintTestsForFunction(function GoFunction, tests []Test) (res string) {
	res = fmt.Sprintf("func Test_%s(t *testing.T) {\n", function.Name())
	for i, test := range tests {
		res += PrintOneTest(i, function, test)
	}
	res += "}\n\n"
	return
}

func PrintTestsForFile(file GoFile, testSet TestSet) (res string) {
	res = ""
	functions := file.Functions()
	for _, function := range functions {
		if tests, ok := testSet[function.Name()]; ok {
			res += PrintTestsForFunction(function, tests)
		}
	}
	return
}

func PrintTestingFile(file GoFile, testSet TestSet) string {
	return PrintPackageDeclaration(file) +
			PrintImportedLibraries() +
			PrintTestsForFile(file, testSet)
}
