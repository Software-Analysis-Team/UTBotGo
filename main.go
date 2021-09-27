package main

import (
	"io/ioutil"
	"os"
	"regexp"
)

func CreateTestingCode(testedFileName string) string {
	testedFile, _ := NewFile(testedFileName)
	testSet := testedFile.GenerateTests()
	return PrintTestingFile(testedFile, testSet)
}

func CreateTestingFile(testedFileName string) {
	isGo, _ := regexp.MatchString(`.*\.go$`, testedFileName)
	isGoTest, _ := regexp.MatchString(`.*_test\.go$`, testedFileName)
	if !isGo || isGoTest {
		return
	}
	testingFileName := testedFileName[:len(testedFileName)-3] + "_test.go"
	testingFile, _ := os.OpenFile(
		testingFileName,
		os.O_WRONLY|os.O_CREATE,
		0644)
	testingCode := CreateTestingCode(testedFileName)
	_, _ = testingFile.WriteString(testingCode)
	_ = testingFile.Close()
}

func CreateTestingDir(testedDirName string, recursively bool) {
	objs, _ := ioutil.ReadDir(testedDirName)
	var CreateTestingSomething func(string)
	if recursively {
		CreateTestingSomething = func(name string) { CreateTestingObject(name, recursively) }
	} else {
		CreateTestingSomething = CreateTestingFile
	}
	for _, obj := range objs {
		CreateTestingSomething(testedDirName + obj.Name())
	}
}

func CreateTestingObject(testedObjectName string, recursively bool) {
	info, _ := os.Stat(testedObjectName)
	switch mode := info.Mode(); {
	case mode.IsDir():
		CreateTestingDir(testedObjectName, recursively)
	case mode.IsRegular():
		CreateTestingFile(testedObjectName)
	}
}

func main() {
	CreateTestingObject(os.Args[1], false)
}
