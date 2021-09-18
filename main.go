package main

import (
    "strings"
    "fmt"
    "go/parser"
    "go/token"
    "go/ast"
    "math/rand"
    "strconv"
    "os"
    "io/ioutil"
    "regexp"
)

type FuncTest []string

type FuncTests []FuncTest

type FileTest map[string]FuncTests

func GenerateString (length int) (str string) {
    n := rand.Intn(length)
    str = ""
    for i := 0; i < n; i++ {
        str += string(rand.Intn(128))
    }
    return
}

func GenerateFuncTests (testedFunc *ast.FuncDecl, funcTestAmt int) (funcTests FuncTests) {
    funcTests = make([]FuncTest, funcTestAmt)
    paramTypes := []string{}
    for _, param := range testedFunc.Type.Params.List {
        typeName := param.Type.(*ast.Ident).Name
        for range param.Names {
            paramTypes = append(paramTypes, typeName)
        }
    }
    paramAmt := len(paramTypes)
    for i := range funcTests {
        funcTests[i] = make(FuncTest, paramAmt)
        for j := range paramTypes {
            switch paramTypes[j] {
            case "int":
                funcTests[i][j] = strconv.Itoa(int(rand.Uint64()))
            case "string":
                funcTests[i][j] = fmt.Sprintf("%#v", GenerateString(42))
            default:
                funcTests[i][j] = "nil"
            }
        }
    }
    return
}

func GenerateFileTest (testedFile string, funcTestAmt int) (fileTest FileTest) {
    fset := token.NewFileSet()
    file, _ := parser.ParseFile(fset, testedFile, nil, 0)
    fileTest = make(FileTest)
    for _, d := range file.Decls {
        if fun, ok := d.(*ast.FuncDecl); ok {
            fileTest[fun.Name.Name] = GenerateFuncTests(fun, funcTestAmt)
        }
    }
    return
}

func GenerateTestingCode (packageName string, fileTest FileTest) (testingCode string) {
    testingCode = "package " + packageName + "\n\nimport \"testing\""
    for name, tests := range fileTest {
        testingCode += "\n\nfunc Test_" + name + " (t *testing.T) {\n"
        for i, test := range tests {
            testStr := strings.Join(test, ", ")
            funcStr := name + "(" + testStr + ")"
            iStr := strconv.Itoa(i)
            testingCode +=
                "    t.Run(\"" + iStr + "\", func (t *testing.T) { defer func () { if recover() != nil { t.Fail() } }(); " + funcStr + " })\n"
        }
        testingCode += "}"
    }
    return
}

func TakePackageName (fileName string) string {
    fset := token.NewFileSet()
    file, _ := parser.ParseFile(fset, fileName, nil, parser.PackageClauseOnly)
    return file.Name.Name
}

func main () {
    dirName := "test-files/"
    files, _ := ioutil.ReadDir(dirName)
    for _, file := range files {
        fileName := file.Name()
        match_go, _ := regexp.MatchString(`.*\.go$`, fileName)
        match_test, _ := regexp.MatchString(`.*_test\.go$`, fileName)
        if match_go && ! match_test {
            testedFile := dirName + fileName
            packageName := TakePackageName (testedFile)
            fileTest := GenerateFileTest(testedFile, 5)
            testingCode := GenerateTestingCode(packageName, fileTest)
            testingFile, _ := os.OpenFile(dirName + fileName[:len(fileName)-3] + "_test.go", os.O_WRONLY|os.O_CREATE, 0644)
            defer testingFile.Close()
            testingFile.WriteString(testingCode)
        }
    }
}
