package main

import (
	"strconv"
	"strings"
)

func PrintPackageDeclaration(file File) string {
	return "package " + file.PackageName() + "\n\n"
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

func PrintCallingOfTestedFunction(function Function, test Test) string {
	return function.Name() + "(" + PrintDataFromTest(test) + ")"
}

func PrintOneTest(index int, function Function, test Test) string {
	return "    t.Run(" +
		"\"" + strconv.Itoa(index+1) + "\", " +
		"func (t *testing.T) { " +
		"defer func () { if recover() != nil { t.Fail() } }(); " +
		PrintCallingOfTestedFunction(function, test) +
		" }" +
		")\n"
}

func PrintTestsForFunction(function Function, tests []Test) (res string) {
	res = "func Test_" + function.Name() + " (t *testing.T) {\n"
	for i, test := range tests {
		res += PrintOneTest(i, function, test)
	}
	res += "}\n\n"
	return
}

func PrintTestsForFile(file File, testSet TestSet) (res string) {
	res = ""
	functions := file.Functions()
	for _, function := range functions {
		if tests, ok := testSet[function.Name()]; ok {
			res += PrintTestsForFunction(function, tests)
		}
	}
	return
}

func PrintTestingFile(file File, testSet TestSet) string {
	return PrintPackageDeclaration(file) +
		PrintImportedLibraries() +
		PrintTestsForFile(file, testSet)
}
